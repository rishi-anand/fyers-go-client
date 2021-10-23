package fyerswatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/rishi-anand/fyers-go-client/api"

	"github.com/rishi-anand/fyers-go-client/utils"

	"github.com/lunixbochs/struc"

	log "github.com/sirupsen/logrus"

	"github.com/sacOO7/gowebsocket"
)

const (
	notifierUrl = "wss://api.fyers.in/socket/v2/dataSock?access_token=%s:%s"
	dataApi     = "https://api.fyers.in/data-rest/v2/quotes/?symbols=%s"
)
const (
	fyPLenHeader      = 24
	fyPLenComnPayload = 48
	fyPLenExtra7208   = 32
	fyPLenBidAsk      = 12
)

type watchNotifier struct {
	conn *gowebsocket.Socket
	nt   api.NotificationType

	tokenMap map[string]string

	apiKey      string
	accessToken string

	onMessage     func(api.Notification)
	onNoReconnect func(int)
	onReconnect   func(int, time.Duration)
	onConnect     func()
	onClose       func()
	onError       func(error)
}

func NewNotifier(apiKey, accessToken string) *watchNotifier {
	return &watchNotifier{
		apiKey:      apiKey,
		accessToken: accessToken,
		tokenMap:    make(map[string]string),
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

//type fyersTickerInput struct {
//	Type    string   `json:"T"`
//	Symbols []string `json:"L2LIST"`
//	SubT    int      `json:"SUB_T"`
//}

func (w *watchNotifier) Subscribe(nt api.NotificationType, symbols ...string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(fmt.Sprintf(notifierUrl, w.apiKey, w.accessToken))

	socket.OnConnectError = w.OnConnectError
	socket.OnTextMessage = w.OnTextMessage
	socket.OnPingReceived = w.OnPingReceived
	socket.OnPongReceived = w.OnPongReceived
	socket.OnDisconnected = w.OnDisconnected
	socket.OnConnected = func(socket gowebsocket.Socket) {
		w.onConnected(socket, symbols)
	}
	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		w.OnBinaryMessage(nt, data, socket)
	}
	socket.Connect()
	w.conn = &socket
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}

func (w *watchNotifier) onConnected(socket gowebsocket.Socket, symbols []string) {
	log.Println("Connected to server")
	socket.SendBinary([]byte(`{"T": "SUB_L2", "L2LIST": [` + strings.Join(symbols, ",") + `], "SUB_T": 1}`))
	if w.onConnect != nil {
		w.onConnect()
	}
}

func (w *watchNotifier) OnConnectError(err error, socket gowebsocket.Socket) {
	socket.Close()
	if w.onError != nil {
		w.onError(err)
	}
}

func (w *watchNotifier) OnTextMessage(message string, socket gowebsocket.Socket) {
	log.Println("Recieved message " + message)
}

func (w *watchNotifier) OnPingReceived(data string, socket gowebsocket.Socket) {
	log.Debugln("Recieved ping " + data)
}

func (w *watchNotifier) OnPongReceived(data string, socket gowebsocket.Socket) {
	log.Debugln("Recieved pong " + data)
}

func (w *watchNotifier) OnDisconnected(err error, socket gowebsocket.Socket) {
	log.Println("Disconnected from server ")
	//reconnect
}

func (w *watchNotifier) OnBinaryMessage(nt api.NotificationType, data []byte, socket gowebsocket.Socket) {
	n := api.Notification{Type: nt}
	if nt == api.SymbolDataTick {
		v := bytes.NewReader(data[0:fyPLenHeader])
		header := &PacketHeader{}
		err := struc.Unpack(v, header)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(header)

		x := bytes.NewReader(data[fyPLenHeader:])
		msg := &PacketMsg{}
		err = struc.Unpack(x, msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg)

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
			err = struc.Unpack(y, extraMsg)
			if err != nil {
				fmt.Println(err)
				return
			}
			n.SymbolData.LastTradedQty = int(extraMsg.Ltq)
			n.SymbolData.LastTradedTime = utils.ToIstTimeFromEpoch(int64(extraMsg.Ltt))
			n.SymbolData.AvgTradedPrice = float32(extraMsg.Atp)
			n.SymbolData.VolumeTradedToday = int64(extraMsg.Vtt)
			n.SymbolData.TotalBuyQty = int64(extraMsg.TotBuy)
			n.SymbolData.TotalSellQty = int64(extraMsg.TotSell)

			depth := make([]api.MarketDepth, 0, 1)
			//market depth to be run 10 times
			msg := data[fyPLenHeader:][fyPLenComnPayload:][fyPLenExtra7208:]
			for i := 0; i < 10; i++ {
				z := bytes.NewReader(msg[:fyPLenBidAsk])
				bidAsk := &PacketMsgMarketDepth{}
				err = struc.Unpack(z, bidAsk)
				if err != nil {
					fmt.Println(err)
					return
				}
				depth = append(depth, api.MarketDepth{Price: float32(bidAsk.Price), Qty: int64(bidAsk.Qty), NumOfOrders: int64(bidAsk.NumOrd)})
			}
			n.SymbolData.MarketPic = depth
		}
	}

	if w.onMessage != nil {
		w.onMessage(n)
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

/*
func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	url := "wss://api.fyers.in/socket/v2/dataSock?access_token=58LA0JLQLD-100:eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJhcGkuZnllcnMuaW4iLCJpYXQiOjE2MzQ4ODM4NjIsImV4cCI6MTYzNDk0OTAyMiwibmJmIjoxNjM0ODgzODYyLCJhdWQiOlsieDoyIiwiZDoyIiwiZDoxIiwieDoxIiwieDowIl0sInN1YiI6ImFjY2Vzc190b2tlbiIsImF0X2hhc2giOiJnQUFBQUFCaGNsa1dZbVMweE9QeWFQRXJlNnRzaEJUX3YtcXc2NFhLUTdpY1RjUFZEdzRGMm1jbU1ZMlNCX0pzYVcxNHhUU2pZTHctZ2kzR18yNE9naW1BWnY2THFpZmhxeEQtcVBFN3Q2NTZreDVBZ3lfazJmST0iLCJkaXNwbGF5X25hbWUiOiJSSVNISSBBTkFORCIsImZ5X2lkIjoiWFIwNzgzOCIsImFwcFR5cGUiOjEwMCwicG9hX2ZsYWciOiJOIn0.7x3bhoTe1MSV3sOml6mqrwzc6H0v581wYwINc44O_IU"
	socket := gowebsocket.New(url)

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
		symbol := `{"T": "SUB_L2", "L2LIST": ["NSE:BANKNIFTY21OCT38200PE"], "SUB_T": 1}`
		socket.SendBinary([]byte(symbol))
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	}

	// https://jsoneditoronline.org/#left=cloud.9dab8a41f431414eb33b376bc501f43e&right=local.vinaja
	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		//fmt.Println(string(data))
		z := bytes.NewReader(data[0:24])
		o := &PacketHeader{}
		err := struc.Unpack(z, o)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(o)

		x := bytes.NewReader(data[24:])
		m := &PacketMsg{}
		err = struc.Unpack(x, m)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(m)

		y := bytes.NewReader(data[24:][48:])
		mk := &PacketMsgExtra{}
		err = struc.Unpack(y, mk)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("=======", mk)

		//market depth to be run 10 times
		rawData := data[24:][48:][32:]
		w := bytes.NewReader(rawData[:12])
		nk := &PacketMsgMarketDepth{}
		err = struc.Unpack(w, nk)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("======= MARKET ", nk)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		//reconnect
		return
	}

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}
*/
/*
Type	Size
uint	Platform Dependent
uint8	8 bits/1 byte
uint16	16 bits/2 byte
uint32	32 bits/4 byte
uint64	64 bits/8 byte

https://docs.python.org/3/library/struct.html

Struct("> 10I Q")
*/
