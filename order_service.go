package main

import (
	"encoding/json"
	"fmt"

	"github.com/rishi-anand/fyers-go-client/api"
	"github.com/rishi-anand/fyers-go-client/utils"
)

func (c *Client) PlaceOrder(order api.Order) api.OrderResponse {
	if resp, err := c.invoke(utils.POST, SingleOrderUrl, order); err != nil {
		return api.OrderResponse{IsFailed: true, Message: err.Error()}
	} else {
		if utils.IsSuccessResponse(resp) {
			var response api.OrderResponse
			if err := json.Unmarshal(resp, &response); err != nil {
				return api.OrderResponse{IsFailed: true, Message: fmt.Errorf("failed to unmarshall order response. %v", err).Error()}
			}
			return response
		} else {
			return api.OrderResponse{IsFailed: true, Message: utils.GetJsonValueAtPath(resp, "message")}
		}
	}
}

func (c *Client) PlaceOrders(orders []api.Order) api.MultiOrderResponse {
	if resp, err := c.invoke(utils.POST, MultiOrderUrl, orders); err != nil {
		return api.MultiOrderResponse{IsFailed: true, Message: err.Error()}
	} else {
		if utils.IsSuccessResponse(resp) {
			var response api.MultiOrderResponse
			if err := json.Unmarshal(resp, &response); err != nil {
				return api.MultiOrderResponse{IsFailed: true, Message: fmt.Errorf("failed to unmarshall order response. %v", err).Error()}
			}
			return response
		} else {
			return api.MultiOrderResponse{IsFailed: true, Message: utils.GetJsonValueAtPath(resp, "message")}
		}
	}
}
