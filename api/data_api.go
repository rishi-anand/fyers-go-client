package api

import (
	"time"

	"github.com/rishi-anand/fyers-go-client/utils"
)

// DataQuote is the response api model for quotes
type DataQuote struct {
	Symbol             string  `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	ShortName          string  `json:"short_name,omitempty" yaml:"short_name,omitempty"`
	Exchange           string  `json:"exchange,omitempty" yaml:"exchange,omitempty"`
	Description        string  `json:"description,omitempty" yaml:"description,omitempty"`
	OriginalName       string  `json:"original_name,omitempty" yaml:"original_name,omitempty"`
	FyToken            string  `json:"fyToken,omitempty" yaml:"fyToken,omitempty"`
	Volume             int64   `json:"volume,omitempty" yaml:"volume,omitempty"`
	PreviousClosePrice float32 `json:"prev_close_price,omitempty" yaml:"prev_close_price,omitempty"`
	LowPrice           float32 `json:"low_price,omitempty" yaml:"low_price,omitempty"`
	HighPrice          float32 `json:"high_price,omitempty" yaml:"high_price,omitempty"`
	OpenPrice          float32 `json:"open_price,omitempty" yaml:"open_price,omitempty"`
	Bid                float32 `json:"bid,omitempty" yaml:"bid,omitempty"`
	Ask                float32 `json:"ask,omitempty" yaml:"ask,omitempty"`
	Spread             float32 `json:"spread,omitempty" yaml:"spread,omitempty"`
	LastPrice          float32 `json:"lp,omitempty" yaml:"lp,omitempty"`
	Chp                float32 `json:"chp,omitempty" yaml:"chp,omitempty"`
	Ch                 float32 `json:"ch,omitempty" yaml:"ch,omitempty"`
	Time               int64   `json:"tt,omitempty" yaml:"tt,omitempty"`
}

func (d *DataQuote) IstTimestamp() time.Time {
	if d.Time != 0 {
		return utils.ToIstTimeFromEpoch(d.Time)
	}
	return time.Time{}
}

// HistoricalData is the response api model for historical data with candles model
type HistoricalData struct {
	Status  string      `json:"s,omitempty" yaml:"s,omitempty"`
	Candles [][]float64 `json:"candles,omitempty" yaml:"candles,omitempty"`
}

// MarketDepth is the response api model for data api for market depth api
type MarketDepth struct {
	TotalBuyQty     int             `json:"totalbuyqty,omitempty" yaml:"totalbuyqty,omitempty"`
	TotalSellQty    int             `json:"totalsellqty,omitempty" yaml:"totalsellqty,omitempty"`
	Bids            []MarketDataBid `json:"bids,omitempty" yaml:"bids,omitempty"`
	Asks            []MarketDataBid `json:"ask,omitempty" yaml:"ask,omitempty"`
	Open            float32         `json:"o,omitempty" yaml:"o,omitempty"`
	High            float32         `json:"h,omitempty" yaml:"h,omitempty"`
	Low             float32         `json:"l,omitempty" yaml:"l,omitempty"`
	Close           float32         `json:"c,omitempty" yaml:"c,omitempty"`
	Chp             float32         `json:"chp,omitempty" yaml:"chp,omitempty"`
	Ch              float32         `json:"ch,omitempty" yaml:"ch,omitempty"`
	LastTradedQty   int             `json:"ltq,omitempty" yaml:"ltq,omitempty"`
	LastTradedTime  int64           `json:"ltt,omitempty" yaml:"ltt,omitempty"`
	LastTradedPrice float32         `json:"ltp,omitempty" yaml:"ltp,omitempty"`
	Volume          int             `json:"v,omitempty" yaml:"v,omitempty"`
	Atp             float32         `json:"atp,omitempty" yaml:"atp,omitempty"`
	LowerCkt        float32         `json:"lower_ckt,omitempty" yaml:"lower_ckt,omitempty"`
	UpperCkt        float32         `json:"upper_ckt,omitempty" yaml:"upper_ckt,omitempty"`
	Expiry          string          `json:"expiry,omitempty" yaml:"expiry,omitempty"`
}

// MarketDataBid denotes each bid and ask
type MarketDataBid struct {
	Price  float32 `json:"price,omitempty" yaml:"price,omitempty"`
	Volume int64   `json:"volume,omitempty" yaml:"volume,omitempty"`
	Ord    int     `json:"ord,omitempty" yaml:"ord,omitempty"`
}
