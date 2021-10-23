package api

import "time"

type NotificationType string

const (
	SymbolDataTick  NotificationType = "symbolData"
	OrderUpdateTick NotificationType = "orderUpdate"
)

type Notification struct {
	Type       NotificationType
	SymbolData SymbolDataNotification
}

type SymbolDataNotification struct {
	Symbol            string        `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	Timestamp         time.Time     `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	FyCode            int           `json:"fyCode,omitempty" yaml:"fyCode,omitempty"`
	FyFlag            int           `json:"fyFlag,omitempty" yaml:"fyFlag,omitempty"`
	PktLength         int           `json:"pktLen,omitempty" yaml:"pktLen,omitempty"`
	Ltp               float32       `json:"ltp,omitempty" yaml:"ltp,omitempty"`
	OpenPrice         float32       `json:"open_price,omitempty" yaml:"open_price,omitempty"`
	HighPrice         float32       `json:"high_price,omitempty" yaml:"high_price,omitempty"`
	LowPrice          float32       `json:"low_price,omitempty" yaml:"low_price,omitempty"`
	ClosePrice        float32       `json:"close_price,omitempty" yaml:"close_price,omitempty"`
	MinOpenPrice      float32       `json:"min_open_price,omitempty" yaml:"min_open_price,omitempty"`
	MinHighPrice      float32       `json:"min_high_price,omitempty" yaml:"min_high_price,omitempty"`
	MinLowPrice       float32       `json:"min_low_price,omitempty" yaml:"min_low_price,omitempty"`
	MinClosePrice     float32       `json:"min_close_price,omitempty" yaml:"min_close_price,omitempty"`
	MinVolume         int64         `json:"min_volume,omitempty" yaml:"min_volume,omitempty"`
	LastTradedQty     int           `json:"last_traded_qty,omitempty" yaml:"last_traded_qty,omitempty"`
	LastTradedTime    time.Time     `json:"last_traded_time,omitempty" yaml:"last_traded_time,omitempty"`
	AvgTradedPrice    float32       `json:"avg_trade_price,omitempty" yaml:"avg_trade_price,omitempty"`
	VolumeTradedToday int64         `json:"vol_traded_today,omitempty" yaml:"vol_traded_today,omitempty"`
	TotalBuyQty       int64         `json:"tot_buy_qty,omitempty" yaml:"tot_buy_qty,omitempty"`
	TotalSellQty      int64         `json:"tot_sell_qty,omitempty" yaml:"tot_sell_qty,omitempty"`
	MarketPic         []MarketDepth `json:"market_pic,omitempty" yaml:"market_pic,omitempty"`
}

type MarketDepth struct {
	Price       float32 `json:"price,omitempty" yaml:"price,omitempty"`
	Qty         int64   `json:"qty,omitempty" yaml:"qty,omitempty"`
	NumOfOrders int64   `json:"num_orders,omitempty" yaml:"num_orders,omitempty"`
}
