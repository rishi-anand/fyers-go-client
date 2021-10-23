package api

import (
	"time"

	"github.com/rishi-anand/fyers-go-client/utils"
)

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
