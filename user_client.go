package fyers

import (
	"encoding/json"
	"fmt"

	"github.com/rishi-anand/fyers-go-client/api"
	"github.com/rishi-anand/fyers-go-client/utils"
)

func (c *client) GetProfile() (api.UserProfile, error) {
	if resp, err := c.invoke(utils.GET, c.toUri(ApiV2, ProfileUrl), nil); err != nil {
		return api.UserProfile{}, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response api.UserProfile
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(resp, "data")), &response); err != nil {
				return api.UserProfile{}, err
			} else {
				return response, nil
			}
		} else {
			return api.UserProfile{}, fmt.Errorf("failed to get profile information. %v", utils.GetJsonValueAtPath(resp, "message"))
		}
	}
}

func (c *client) GetFund() (api.UserFund, error) {
	if resp, err := c.invoke(utils.GET, c.toUri(ApiV2, FundsUrl), nil); err != nil {
		return api.UserFund{}, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response api.UserFund
			if json.Unmarshal(resp, &response); err != nil {
				return api.UserFund{}, err
			} else {
				return response, nil
			}
		} else {
			return api.UserFund{}, fmt.Errorf("failed to get fund information. %v", utils.GetJsonValueAtPath(resp, "message"))
		}
	}
}

func (c *client) GetHoldings() (api.UserHoldings, error) {
	if resp, err := c.invoke(utils.GET, c.toUri(ApiV2, HoldingsUrl), nil); err != nil {
		return api.UserHoldings{}, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response api.UserHoldings
			if json.Unmarshal(resp, &response); err != nil {
				return api.UserHoldings{}, err
			} else {
				return response, nil
			}
		} else {
			return api.UserHoldings{}, fmt.Errorf("failed to get holdings information. %v", utils.GetJsonValueAtPath(resp, "message"))
		}
	}
}

func (c *client) GetPositions() (api.UserPosition, error) {
	if resp, err := c.invoke(utils.GET, c.toUri(ApiV2, PositionsUrl), nil); err != nil {
		return api.UserPosition{}, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response api.UserPosition
			if json.Unmarshal(resp, &response); err != nil {
				return api.UserPosition{}, err
			} else {
				return response, nil
			}
		} else {
			return api.UserPosition{}, fmt.Errorf("failed to get positions information. %v", utils.GetJsonValueAtPath(resp, "message"))
		}
	}
}

func (c *client) ListTrades() ([]api.TradeBook, error) {
	if resp, err := c.invoke(utils.GET, c.toUri(ApiV2, TradeBookUrl), nil); err != nil {
		return nil, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response []api.TradeBook
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(resp, "tradeBook")), &response); err != nil {
				return nil, err
			} else {
				return response, nil
			}
		} else {
			return nil, fmt.Errorf("failed to get trades information. %v", utils.GetJsonValueAtPath(resp, "message"))
		}
	}
}
