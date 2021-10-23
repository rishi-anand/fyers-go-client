package api

import "time"

type OrderType int

const (
	LimitOrder     OrderType = 1
	MarketOrder    OrderType = 2
	StopLossOrder  OrderType = 3
	StopLimitOrder OrderType = 4
)

type ProductType string

const (
	CNCOrder      ProductType = "CNC"
	IntradayOrder ProductType = "INTRADAY"
	MarginOrder   ProductType = "MARGIN"
	CoverOrder    ProductType = "CO"
	BracketOrder  ProductType = "BO"
)

type OrderStatus int

const (
	Cancelled    OrderStatus = 1
	Traded       OrderStatus = 2
	ForFutureUse OrderStatus = 3
	Transit      OrderStatus = 4
	Rejected     OrderStatus = 5
	Pending      OrderStatus = 6
	Expired      OrderStatus = 7
)

type OrderSide int

const (
	BuyOrder  OrderSide = 1
	SellOrder OrderSide = -1
)

type ValidityType string

const (
	ImmediateOrCancelValidity ValidityType = "IOC"
	EndOfDayValidity          ValidityType = "DAY"
)

type OrderTime string

const (
	MarketOpenOrder  OrderTime = "False"
	AfterMarketOrder OrderTime = "True"
)

// More information can be found at https://api-docs.fyers.in/v2/#single-order
type Order struct {
	// Eg: NSE:SBIN-EQ
	Symbol string `json:"symbol" yaml:"symbol"`
	// The quantity should be in multiples of lot size for derivatives.
	Qty int `json:"qty" yaml:"qty"`
	/*
		Possible Values	Description
		1	Limit order
		2	Market order
		3	Stop order (SL-M)
		4	Stoplimit order (SL-L)
	*/
	Type OrderType `json:"type" yaml:"type"`
	/*
		Possible Values	Description
		1	Buy
		-1	Sell
	*/
	OrderSide OrderSide `json:"side" yaml:"side"`
	/*
		CNC => For equity only
		INTRADAY => Applicable for all segments.
		MARGIN => Applicable only for derivatives
		CO => Cover Order
		BO => Bracket Order
	*/
	ProductType ProductType `json:"productType" yaml:"productType"`
	/*
		Default => 0
		Provide valid price for Limit and Stoplimit orders
	*/
	LimitPrice float32 `json:"limitPrice" yaml:"limitPrice"`
	/*
		Default => 0
		Provide valid price for Stop and Stoplimit orders
	*/
	StopPrice float32 `json:"stopPrice" yaml:"stopPrice"`
	/*
		IOC => Immediate or Cancel
		DAY => Valid till the end of the day
	*/
	Validity ValidityType `json:"validity" yaml:"validity"`
	/*
		Default => 0
		Allowed only for Equity
	*/
	DisclosedQty int `json:"disclosedQty" yaml:"disclosedQty"`
	/*
		False => When market is open
		True => When placing AMO order
	*/
	OrderTime OrderTime `json:"offlineOrder" yaml:"offlineOrder"`
	// Provide valid price for CO and BO orders
	StopLoss float32 `json:"stopLoss" yaml:"stopLoss"`
	/*
		Default => 0
		Provide valid price for BO orders
	*/
	TakeProfit float32 `json:"takeProfit" yaml:"takeProfit"`
}

type OrderOperation string

const (
	CreateOrder    OrderOperation = "CREATE"
	ModifyOrder    OrderOperation = "MODIFY"
	CancelOrder    OrderOperation = "CANCEL"
	SquareOffOrder OrderOperation = "SQUAREOFF"
)

type OrderQuantity string

const (
	SingleQuantity OrderQuantity = "Single"
	MultiOrder     OrderQuantity = "Multi"
)

type OrderId struct {
	Id string `json:"id" yaml:"id"`
}

func NewOrderId(id string) OrderId {
	return OrderId{Id: id}
}

type OrderUpdate struct {
	Id string `json:"id,omitempty" yaml:"id,omitempty"`
	// Optional. Only incase of Limit/ Stoplimit orders
	LimitPrice float32 `json:"limitPrice,omitempty" yaml:"limitPrice,omitempty"`
	// Optional. Only incase of Stop/ Stoplimit orders
	StopLoss float32 `json:"takeProfit,omitempty" yaml:"takeProfit,omitempty"`
	// Optional. Incase you want to modify the quantity
	Qty int `json:"qty,omitempty" yaml:"qty,omitempty"`
	/*
		Possible Values	Description
		1	Limit order
		2	Market order
		3	Stop order (SL-M)
		4	Stoplimit order (SL-L)
	*/
	Type OrderType `json:"type,omitempty" yaml:"type,omitempty"`
}

type OrderResponse struct {
	Id      string `json:"id,omitempty" yaml:"id,omitempty"`
	Status  string `json:"s,omitempty" yaml:"s,omitempty"`
	Code    int    `json:"code,omitempty" code:"symbol,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	IsFailed bool `json:"isFailed,omitempty" yaml:"isFailed,omitempty"` //custom
}

type MultiOrderResponse struct {
	Data    []OrderDataResponse `json:"data,omitempty" yaml:"data,omitempty"`
	Status  string              `json:"s,omitempty" yaml:"s,omitempty"`
	Code    int                 `json:"code,omitempty" code:"symbol,omitempty"`
	Message string              `json:"message,omitempty" yaml:"message,omitempty"`

	IsFailed bool `json:"isFailed,omitempty" yaml:"isFailed,omitempty"` //custom
}

type OrderDataResponse struct {
	StatusCode        int           `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
	Body              OrderResponse `json:"body,omitempty" yaml:"body,omitempty"`
	StatusDescription string        `json:"statusDescription,omitempty" yaml:"statusDescription,omitempty"`
}

type OrderBookResponse struct {
	Status  string `json:"s,omitempty" yaml:"s,omitempty"`
	Code    int    `json:"code,omitempty" code:"symbol,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	IsFailed bool `json:"isFailed,omitempty" yaml:"isFailed,omitempty"` //custom

	OrderBook []OrderBook `json:"orderBook,omitempty" yaml:"orderBook,omitempty"`
}

type OrderBook struct {
	Id             string  `json:"id,omitempty" yaml:"id,omitempty"`
	OrderTime      string  `json:"orderDateTime,omitempty" yaml:"orderDateTime,omitempty"`
	Side           int     `json:"side,omitempty" yaml:"side,omitempty"`
	Segment        int     `json:"segment,omitempty" yaml:"segment,omitempty"`
	ProductType    string  `json:"productType,omitempty" yaml:"productType,omitempty"`
	Status         int     `json:"status,omitempty" yaml:"status,omitempty"`
	Quantity       int     `json:"qty,omitempty" yaml:"qty,omitempty"`
	LimitPrice     float32 `json:"limitPrice,omitempty" yaml:"limitPrice,omitempty"`
	StopPrice      float32 `json:"stopPrice,omitempty" yaml:"stopPrice,omitempty"`
	Type           int     `json:"type,omitempty" yaml:"type,omitempty"`
	Validity       string  `json:"orderValidity,omitempty" yaml:"orderValidity,omitempty"`
	Source         string  `json:"Source,omitempty" yaml:"Source,omitempty"`
	Exchange       int     `json:"exchange,omitempty" yaml:"exchange,omitempty"`
	SerialNo       int     `json:"slNo,omitempty" yaml:"slNo,omitempty"`
	OfflineOrder   bool    `json:"offlineOrder,omitempty" yaml:"offlineOrder,omitempty"`
	Message        string  `json:"message,omitempty" yaml:"message,omitempty"`
	TradedPrice    float32 `json:"tradedPrice,omitempty" yaml:"tradedPrice,omitempty"`
	Symbol         string  `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	ExSymbol       string  `json:"ex_sym,omitempty" yaml:"ex_sym,omitempty"`
	OrderNumStatus string  `json:"orderNumStatus,omitempty" yaml:"orderNumStatus,omitempty"`
	Description    string  `json:"description,omitempty" yaml:"description,omitempty"`
}

func (f OrderBook) GetStatus() OrderStatus {
	return OrderStatus(f.Status)
}

func (f OrderBook) GetOrderTime() time.Time {
	// 22-Aug-2021 19:26:29
	orderTime, _ := time.Parse("02-Jan-2006 15:04:05", f.OrderTime)
	return orderTime
}

func (f OrderBook) GetSide() OrderSide {
	return OrderSide(f.Side)
}

func (f OrderBook) GetProductType() ProductType {
	return ProductType(f.ProductType)
}

func (f OrderBook) GetOrderType() OrderType {
	return OrderType(f.Type)
}

func (f OrderBook) IsOfflineOrder() OrderTime {
	if f.OfflineOrder {
		return AfterMarketOrder
	} else {
		return MarketOpenOrder
	}
}
