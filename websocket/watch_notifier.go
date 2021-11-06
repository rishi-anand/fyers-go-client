package fyerswatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/rishi-anand/fyers-go-client/api"

	"github.com/rishi-anand/fyers-go-client/utils"

	"github.com/lunixbochs/struc"

	log "github.com/sirupsen/logrus"

	"github.com/sacOO7/gowebsocket"
)

const (
	notifierUrl      = "wss://api.fyers.in/socket/v2/dataSock?access_token=%s:%s"
	orderNotifierUrl = "wss://api.fyers.in/socket/v2/orderSock?type=orderUpdate&access_token=%s:%s&user-agent=fyers-api"
	dataApi          = "https://api.fyers.in/data-rest/v2/quotes/?symbols=%s"
)

const (
	fyPLenHeader      = 24
	fyPLenComnPayload = 48
	fyPLenExtra7208   = 32
	fyPLenBidAsk      = 12

	maxAttempt = 50
	delay      = 10 * time.Second
)

var subsLatch sync.Mutex

type watchNotifier struct {
	conn gowebsocket.Socket
	nt   api.NotificationType

	failedAttempt int
	closeConnReq  bool

	tokenMap          map[string]string
	subscribedSymbols map[string]bool

	apiKey      string
	accessToken string

	onMessage     func(api.Notification)
	onNoReconnect func(int)
	onReconnect   func(int, time.Duration)
	onConnect     func()
	onClose       func()
	onError       func(error)

	interrupt chan os.Signal
}

func NewNotifier(apiKey, accessToken string) *watchNotifier {
	return &watchNotifier{
		apiKey:            apiKey,
		accessToken:       accessToken,
		tokenMap:          make(map[string]string),
		subscribedSymbols: make(map[string]bool),
	}
}

func (w *watchNotifier) WithOnMessageFunc(f func(api.Notification)) *watchNotifier {
	w.onMessage = f
	return w
}

func (w *watchNotifier) WithOnConnectFunc(f func()) *watchNotifier {
	w.onConnect = f
	return w
}

func (w *watchNotifier) WithOnErrorFunc(f func(err error)) *watchNotifier {
	w.onError = f
	return w
}

func (w *watchNotifier) WithOnCloseFunc(f func()) *watchNotifier {
	w.onClose = f
	return w
}

func (w *watchNotifier) Disconnect() {
	w.closeConnReq = true
	if w.conn.IsConnected {
		w.conn.Close()
	}
}

func (w *watchNotifier) Unsubscribe(symbols ...string) {
	subsLatch.Lock()
	defer subsLatch.Unlock()
	log.Println("Unsubscribing from server")
	if w.conn.IsConnected {
		if w.nt == api.SymbolDataTick {
			if len(symbols) > 0 {
				unSubsL := w.deleteFromSubsList(symbols...)
				if len(unSubsL) > 0 {
					w.conn.SendBinary([]byte(`{"T": "SUB_L2", "L2LIST": [` + strings.Join(utils.FormatStrArrWithQuotes(unSubsL), ",") + `], "SUB_T": 0}`))
				}
			}
		} else if w.nt == api.OrderUpdateTick {
			w.conn.SendBinary([]byte(`{"T": "SUB_ORD", "SLIST": "orderUpdate", "SUB_T": 0}`))
		}
	}
}

func (w *watchNotifier) Subscribe(nt api.NotificationType, symbols ...string) {
	subsLatch.Lock()
	if !w.conn.IsConnected {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		w.nt = nt
		var socket gowebsocket.Socket
		if nt == api.SymbolDataTick {
			w.setFyersTokenForSymbols(symbols)
			socket = gowebsocket.New(fmt.Sprintf(notifierUrl, w.apiKey, w.accessToken))
		} else {
			socket = gowebsocket.New(fmt.Sprintf(orderNotifierUrl, w.apiKey, w.accessToken))
		}

		socket.OnConnectError = w.OnConnectError
		socket.OnTextMessage = w.OnTextMessage
		socket.OnPingReceived = w.OnPingReceived
		socket.OnPongReceived = w.OnPongReceived
		socket.OnDisconnected = w.OnDisconnected
		socket.OnConnected = func(socket gowebsocket.Socket) {
			w.onConnected(socket, nt, w.addToSubsList(symbols...)...)
		}

		socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
			w.OnBinaryMessage(socket, nt, data)
		}
		socket.Connect()
		subsLatch.Unlock()
		w.conn = socket
		for {
			select {
			case <-interrupt:
				socket.Close()
				return
			}
		}
	} else {
		if w.nt != nt {
			log.Errorf("current subscription is for %v notification, but received add subscription for %v notification", w.nt, nt)
			w.onError(fmt.Errorf("current subscription is for %v notification, but received add subscription for %v notification", w.nt, nt))
			return
		}
		w.onConnected(w.conn, nt, w.addToSubsList(symbols...)...)
		subsLatch.Unlock()
	}
}

func (w *watchNotifier) reconnect() {
	if !w.closeConnReq {
		w.conn.Connect()
	}
}

func (w *watchNotifier) onConnected(socket gowebsocket.Socket, nt api.NotificationType, symbols ...string) {
	log.Println("Connected to server")
	if nt == api.SymbolDataTick {
		if len(symbols) > 0 {
			socket.SendBinary([]byte(`{"T": "SUB_L2", "L2LIST": [` + strings.Join(utils.FormatStrArrWithQuotes(symbols), ",") + `], "SUB_T": 1}`))
		}
	} else if nt == api.OrderUpdateTick {
		socket.SendBinary([]byte(`{"T": "SUB_ORD", "SLIST": "orderUpdate", "SUB_T": 1}`))
	}
	if w.onConnect != nil {
		w.onConnect()
	}
}

func (w *watchNotifier) OnConnectError(err error, socket gowebsocket.Socket) {
	socket.Close()
	w.failedAttempt++
	if w.failedAttempt < maxAttempt {
		w.reconnect()
		time.Sleep(delay)
	}
	w.notifyError(err)
}

type WsTextMsg struct {
	Code    int    `json:"code,omitempty" yaml:"code,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Status  string `json:"s,omitempty" yaml:"s,omitempty"`
}

func (w *watchNotifier) OnTextMessage(msg string, socket gowebsocket.Socket) {
	var e WsTextMsg
	log.Println(msg)
	_ = json.Unmarshal([]byte(msg), &e)
	if e.Status == "error" {
		if e.Code == -1600 {
			if socket.IsConnected {
				w.closeConnReq = true
				socket.Close()
				w.interrupt <- os.Interrupt
			}
		}

		if w.onError != nil {
			w.onError(fmt.Errorf("%s", e.Message))
		}
	}
}

func (w *watchNotifier) OnPingReceived(data string, socket gowebsocket.Socket) {
	log.Debugln("Recieved ping " + data)
}

func (w *watchNotifier) OnPongReceived(data string, socket gowebsocket.Socket) {
	log.Debugln("Recieved pong " + data)
}

func (w *watchNotifier) OnDisconnected(err error, socket gowebsocket.Socket) {
	if !w.closeConnReq {
		w.reconnect()
	} else {
		log.Println("Disconnected from server ")
	}
}

func (w *watchNotifier) addToSubsList(symbols ...string) []string {
	newSubsL := make([]string, 0, 1)
	for _, ss := range symbols {
		if _, ok := w.subscribedSymbols[ss]; !ok {
			newSubsL = append(newSubsL, ss)
			w.subscribedSymbols[ss] = true
		}
	}
	return newSubsL
}

func (w *watchNotifier) deleteFromSubsList(symbols ...string) []string {
	unSubsL := make([]string, 0, 1)
	for _, ss := range symbols {
		if _, ok := w.subscribedSymbols[ss]; ok {
			unSubsL = append(unSubsL, ss)
			delete(w.subscribedSymbols, ss)
		}
	}
	return unSubsL
}

func (w *watchNotifier) OnBinaryMessage(socket gowebsocket.Socket, nt api.NotificationType, data []byte) {
	n := api.Notification{Type: nt}
	if nt == api.SymbolDataTick {
		v := bytes.NewReader(data[0:fyPLenHeader])
		header := &PacketHeader{}
		if err := struc.Unpack(v, header); err != nil {
			w.notifyError(err)
			return
		}

		x := bytes.NewReader(data[fyPLenHeader:])
		msg := &PacketMsg{}
		if err := struc.Unpack(x, msg); err != nil {
			w.notifyError(err)
			return
		}

		n.SymbolData = api.SymbolDataNotification{
			Symbol:        w.tokenMap[fmt.Sprintf("%d", header.FyersToken)],
			FyCode:        int(header.FyersCode),
			Timestamp:     utils.ToIstTimeFromEpoch(int64(header.Timestamp)),
			FyFlag:        int(header.Flag),
			PktLength:     int(header.PacketLength),
			Ltp:           float32(msg.Ltp) / float32(msg.Pc),
			OpenPrice:     float32(msg.Op) / float32(msg.Pc),
			HighPrice:     float32(msg.Hp) / float32(msg.Pc),
			LowPrice:      float32(msg.Lp) / float32(msg.Pc),
			ClosePrice:    float32(msg.Cp) / float32(msg.Pc),
			MinOpenPrice:  float32(msg.Mop) / float32(msg.Pc),
			MinHighPrice:  float32(msg.Mhp) / float32(msg.Pc),
			MinLowPrice:   float32(msg.Mlp) / float32(msg.Pc),
			MinClosePrice: float32(msg.Mcp) / float32(msg.Pc),
			MinVolume:     int64(msg.Mv),
		}
		if _, found := fyCodeMap[int(header.FyersCode)]; !found {
			y := bytes.NewReader(data[fyPLenHeader:][fyPLenComnPayload:])
			extraMsg := &PacketMsgExtra{}
			if err := struc.Unpack(y, extraMsg); err != nil {
				w.notifyError(err)
				return
			}
			n.SymbolData.LastTradedQty = int(extraMsg.Ltq)
			n.SymbolData.LastTradedTime = utils.ToIstTimeFromEpoch(int64(extraMsg.Ltt))
			n.SymbolData.AvgTradedPrice = float32(extraMsg.Atp)
			n.SymbolData.VolumeTradedToday = int64(extraMsg.Vtt)
			n.SymbolData.TotalBuyQty = int64(extraMsg.TotBuy)
			n.SymbolData.TotalSellQty = int64(extraMsg.TotSell)

			depth := make([]api.MarketBid, 0, 1)
			//market depth to be run 10 times
			msg := data[fyPLenHeader:][fyPLenComnPayload:][fyPLenExtra7208:]
			for i := 0; i < 10; i++ {
				z := bytes.NewReader(msg[:fyPLenBidAsk])
				bidAsk := &PacketMsgMarketDepth{}
				if err := struc.Unpack(z, bidAsk); err != nil {
					w.notifyError(err)
					return
				}
				depth = append(depth, api.MarketBid{Price: float32(bidAsk.Price), Qty: int64(bidAsk.Qty), NumOfOrders: int64(bidAsk.NumOrd)})
			}
			n.SymbolData.MarketPic = depth
		}
	} else if nt == api.OrderUpdateTick {
		if utils.IsSuccessResponse(data) {
			var order api.OrderNotification
			if err := json.Unmarshal([]byte(utils.GetJsonValueAtPath(data, "d")), &order); err == nil {
				n.OrderData = order
			} else {
				n.OrderData = api.OrderNotification{
					IsError: true,
					Message: utils.GetJsonValueAtPath(data, err.Error()),
				}
			}
		} else {
			n.OrderData = api.OrderNotification{
				IsError: true,
				Message: utils.GetJsonValueAtPath(data, "msg"),
			}
		}
	}

	if w.onMessage != nil {
		w.onMessage(n)
	}
}

func (w *watchNotifier) notifyError(err error) {
	if w.onError != nil {
		w.onError(err)
	}
}

var fyCodeMap = map[int]bool{
	7202: true,
	7207: true,
	27:   true,
}

func (w *watchNotifier) setFyersTokenForSymbols(symbols []string) error {
	headerMap := map[string]string{
		"Authorization": fmt.Sprintf("%s:%s", w.apiKey, w.accessToken),
		"Content-Type":  "application/json",
	}

	if respByte, err := utils.DoHttpCall(utils.GET, fmt.Sprintf(dataApi, strings.Join(symbols, ",")), nil, headerMap); err != nil {
		return err
	} else {
		if utils.IsSuccessResponse(respByte) {
			var quoteResp []api.DataQuote
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(respByte, "d.#.v")), &quoteResp); err != nil {
				return err
			} else {
				for _, q := range quoteResp {
					w.tokenMap[q.FyToken] = q.Symbol
				}
				return nil
			}
		} else {
			return fmt.Errorf("failed to get quote for symbols %v. %v", symbols, utils.GetJsonValueAtPath(respByte, "errmsg"))
		}
	}
}

type PacketHeader struct { // > Q L H H H 6x
	FyersToken   uint64 `struc:"uint64"` //Q | unsigned long long | integer | 8 byte
	Timestamp    uint32 `struc:"uint32"` //L | unsigned long | integer | 4 byte
	FyersCode    uint16 `struc:"uint16"` //H | unsigned short | integer | 2 byte
	Flag         uint16 `struc:"uint16"` //H | unsigned short | integer | 2 byte
	PacketLength uint16 `struc:"uint16"` //H | unsigned short | integer | 2 byte
}

type PacketMsg struct { // > 10I Q
	Pc  uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Ltp uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Op  uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Hp  uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Lp  uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Cp  uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Mop uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Mhp uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Mlp uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Mcp uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Mv  uint64 `struc:"uint64"` //Q | unsigned long long | integer | 8 byte
}

type PacketMsgExtra struct { // > 4I 2Q
	Ltq     uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Ltt     uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Atp     uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Vtt     uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	TotBuy  uint64 `struc:"uint64"` //Q | unsigned long long | integer | 8 byte
	TotSell uint64 `struc:"uint64"` //Q | unsigned long long | integer | 8 byte
}

type PacketMsgMarketDepth struct { // > 3I
	Price  uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	Qty    uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
	NumOrd uint32 `struc:"uint32"` //I | unsigned int | integer | 4 byte
}
