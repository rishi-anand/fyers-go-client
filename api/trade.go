package api

type TradeBook struct {
	Id              string  `json:"id,omitempty" yaml:"id,omitempty"`
	Symbol          string  `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	Row             int     `json:"row,omitempty" yaml:"row,omitempty"`
	ClientId        string  `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	OrderTime       string  `json:"orderDateTime,omitempty" yaml:"orderDateTime,omitempty"`
	OrderNo         string  `json:"orderNumber,omitempty" yaml:"orderNumber,omitempty"`
	ExchangeOrderNo string  `json:"exchangeOrderNo,omitempty" yaml:"exchangeOrderNo,omitempty"`
	Exchange        int     `json:"exchange,omitempty" yaml:"exchange,omitempty"`
	Side            int     `json:"side,omitempty" yaml:"side,omitempty"`
	Segment         int     `json:"segment,omitempty" yaml:"segment,omitempty"`
	OrderType       int     `json:"orderType,omitempty" yaml:"orderType,omitempty"`
	FyToken         string  `json:"fyToken,omitempty" yaml:"fyToken,omitempty"`
	ProductType     string  `json:"productType,omitempty" yaml:"productType,omitempty"`
	TradedQty       int     `json:"tradedQty,omitempty" yaml:"tradedQty,omitempty"`
	TradedPrice     float32 `json:"tradePrice,omitempty" yaml:"tradePrice,omitempty"`
	TradedValue     float64 `json:"tradeValue,omitempty" yaml:"tradeValue,omitempty"`
	TradeNumber     int     `json:"tradeNumber,omitempty" yaml:"tradeNumber,omitempty"`
}
