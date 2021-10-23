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
