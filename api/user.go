package api

type UserProfile struct {
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	Image         string `json:"image,omitempty" yaml:"image,omitempty"`
	DisplayName   string `json:"display_name,omitempty" yaml:"display_name,omitempty"`
	EmailId       string `json:"email_id,omitempty" yaml:"email_id,omitempty"`
	Pan           string `json:"PAN,omitempty" yaml:"PAN,omitempty"`
	FyersId       string `json:"fy_id,omitempty" yaml:"fy_id,omitempty"`
	PwdChangeDate string `json:"pwd_change_date,omitempty" yaml:"pwd_change_date,omitempty"`
	PwdToExpire   int32  `json:"pwd_to_expire,omitempty" yaml:"pwd_to_expire,omitempty"`
}

type UserFund struct {
	FundLimit []Fund `json:"fund_limit,omitempty" yaml:"fund_limit,omitempty"`
}

type Fund struct {
	Id              int     `json:"id,omitempty" yaml:"id,omitempty"`
	Title           string  `json:"title,omitempty" yaml:"title,omitempty"`
	EquityAmount    float64 `json:"equityAmount,omitempty" yaml:"equityAmount,omitempty"`
	CommodityAmount float64 `json:"commodityAmount,omitempty" yaml:"commodityAmount,omitempty"`
}

type UserHoldings struct {
	Holdings    []Holding   `json:"holdings,omitempty" yaml:"holdings,omitempty"`
	HoldingMeta HoldingMeta `json:"overall,omitempty" yaml:"overall,omitempty"`
}

type Holding struct {
	Type         string  `json:"holdingType,omitempty" yaml:"holdingType,omitempty"`
	Symbol       string  `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	Qty          int     `json:"quantity,omitempty" yaml:"quantity,omitempty"`
	CostPrice    float64 `json:"costPrice,omitempty" yaml:"costPrice,omitempty"`
	MarketValue  float64 `json:"marketVal,omitempty" yaml:"marketVal,omitempty"`
	RemainingQty int     `json:"remainingQuantity,omitempty" yaml:"remainingQuantity,omitempty"`
	Pl           float64 `json:"pl,omitempty" yaml:"id,omitempty"`
	Ltp          float64 `json:"ltp,omitempty" yaml:"ltp,omitempty"`
	Id           int     `json:"id,omitempty" yaml:"id,omitempty"`
	FyToken      int64   `json:"fyToken,omitempty" yaml:"fyToken,omitempty"`
	Exchange     int     `json:"exchange,omitempty" yaml:"exchange,omitempty"`
}

type HoldingMeta struct {
	TotalCount        int     `json:"count_total,omitempty" yaml:"count_total,omitempty"`
	TotalInvestment   float64 `json:"total_investment,omitempty" yaml:"total_investment,omitempty"`
	TotalCurrentValue float64 `json:"total_current_value,omitempty" yaml:"total_current_value,omitempty"`
	TotalPL           float64 `json:"total_pl,omitempty" yaml:"total_pl,omitempty"`
	PLPercentage      float32 `json:"pnl_perc,omitempty" yaml:"pnl_perc,omitempty"`
}
