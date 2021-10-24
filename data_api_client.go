package fyers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rishi-anand/fyers-go-client/api"
	"github.com/rishi-anand/fyers-go-client/utils"
)

func (c *client) GetQuote(symbols []string) ([]api.DataQuote, error) {
	if resp, err := c.invoke(utils.GET, fmt.Sprintf(QuoteUrl, strings.Join(symbols, ",")), nil); err != nil {
		return nil, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var quoteResp []api.DataQuote
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(resp, "d.#.v")), &quoteResp); err != nil {
				return nil, err
			} else {
				return quoteResp, nil
			}
		} else {
			return nil, fmt.Errorf("failed to get quote for symbols %v. %v", symbols, utils.GetJsonValueAtPath(resp, "errmsg"))
		}
	}
}

func (c *client) GetHistoricalData(symbol, resolution, start, end string, dateFormat int) (api.HistoricalData, error) {
	queryParam := fmt.Sprintf("?symbol=%s&date_format=%d&cont_flag=1", symbol, dateFormat)
	if len(resolution) > 0 {
		queryParam = fmt.Sprintf("%s&resolution=%s", queryParam, resolution)
	} else {
		queryParam = fmt.Sprintf("%s&resolution=30", queryParam)
	}
	if len(start) > 0 {
		queryParam = fmt.Sprintf("%s&range_from=%s", queryParam, start)
	} else {
		queryParam = fmt.Sprintf("%s&range_from=2021-01-01", queryParam)
	}
	if len(end) > 0 {
		queryParam = fmt.Sprintf("%s&range_to=%s", queryParam, end)
	} else {
		queryParam = fmt.Sprintf("%s&range_to=2021-01-002", queryParam)
	}

	if resp, err := c.invoke(utils.GET, HistoricalDataApiUrl+queryParam, nil); err != nil {
		return api.HistoricalData{}, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response api.HistoricalData
			if json.Unmarshal(resp, &response); err != nil {
				return api.HistoricalData{}, err
			} else {
				return response, nil
			}
		} else {
			return api.HistoricalData{}, fmt.Errorf("failed to get hostorical data for symbol %v. %v", symbol, utils.GetJsonValueAtPath(resp, "errmsg"))
		}
	}
	return api.HistoricalData{}, fmt.Errorf("%s", "NotImplemented")
}

func (c *client) GetMarketDepth(symbol string) (map[string]api.MarketDepth, error) {
	if resp, err := c.invoke(utils.GET, fmt.Sprintf(MarketDepthApiUrl, symbol), nil); err != nil {
		return nil, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response map[string]api.MarketDepth
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(resp, "d")), &response); err != nil {
				return nil, err
			} else {
				return response, nil
			}
		} else {
			return nil, fmt.Errorf("failed to get market depth for symbol %v. %v", symbol, utils.GetJsonValueAtPath(resp, "errmsg"))
		}
	}
}
