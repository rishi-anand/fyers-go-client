package api

type UserPosition struct {
	NetPositions []Position   `json:"netPositions,omitempty" yaml:"netPositions,omitempty"`
	PositionMeta PositionMeta `json:"overall,omitempty" yaml:"overall,omitempty"`
}

type Position struct {
	Symbol           string    `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	Id               string    `json:"id,omitempty" yaml:"id,omitempty"`
	NetQty           int       `json:"netQty,omitempty" yaml:"netQty,omitempty"`
	Qty              int       `json:"qty,omitempty" yaml:"qty,omitempty"`
	AvgPrice         float32   `json:"avgPrice,omitempty" yaml:"avgPrice,omitempty"`
	NetAvg           float32   `json:"netAvg,omitempty" yaml:"netAvg,omitempty"`
	Side             OrderSide `json:"side,omitempty" yaml:"side,omitempty"`
	ProductType      string    `json:"productType,omitempty" yaml:"productType,omitempty"`
	RealizedProfit   float32   `json:"realized_profit,omitempty" yaml:"realized_profit,omitempty"`
	UnrealizedProfit float32   `json:"unrealized_profit,omitempty" yaml:"unrealized_profit,omitempty"`
	Pl               float32   `json:"pl,omitempty" yaml:"pl,omitempty"`
	Ltp              float32   `json:"ltp,omitempty" yaml:"ltp,omitempty"`
	BuyQty           int       `json:"buyQty,omitempty" yaml:"buyQty,omitempty"`
	BuyAvg           float32   `json:"buyAvg,omitempty" yaml:"buyAvg,omitempty"`
	BuyVal           float64   `json:"buyVal,omitempty" yaml:"buyVal,omitempty"`
	SellQty          int       `json:"sellQty,omitempty" yaml:"sellQty,omitempty"`
	SellAvg          float32   `json:"sellAvg,omitempty" yaml:"sellAvg,omitempty"`
	SellVal          float64   `json:"sellVal,omitempty" yaml:"sellVal,omitempty"`
	SlNo             int       `json:"slNo,omitempty" yaml:"slNo,omitempty"`
	FyToken          string    `json:"fyToken,omitempty" yaml:"fyToken,omitempty"`
	CrossCurrency    string    `json:"crossCurrency,omitempty" yaml:"crossCurrency,omitempty"`
	RbiRefRate       int       `json:"rbiRefRate,omitempty" yaml:"rbiRefRate,omitempty"`
	Segment          int       `json:"segment,omitempty" yaml:"segment,omitempty"`
}

type PositionMeta struct {
	TotalCount   int     `json:"count_total,omitempty" yaml:"count_total,omitempty"`
	OpenCount    int     `json:"count_open,omitempty" yaml:"count_open,omitempty"`
	TotalPl      float64 `json:"pl_total,omitempty" yaml:"pl_total,omitempty"`
	RealizedPl   float64 `json:"pl_realized,omitempty" yaml:"pl_realized,omitempty"`
	UnrealizedPl float64 `json:"pl_unrealized,omitempty" yaml:"pl_unrealized,omitempty"`
}
