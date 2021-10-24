package fyers

import (
	"encoding/json"
	"fmt"

	"github.com/rishi-anand/fyers-go-client/api"
	"github.com/rishi-anand/fyers-go-client/utils"
)

func (c *client) PlaceOrder(order api.Order) api.OrderResponse {
	if resp, err := c.invoke(utils.POST, OrdersUrl, order); err != nil {
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

func (c *client) PlaceOrders(orders []api.Order) api.MultiOrderResponse {
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

func (c *client) ListOrders() ([]api.OrderBook, error) {
	if resp, err := c.invoke(utils.GET, OrdersUrl, nil); err != nil {
		return nil, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response []api.OrderBook
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(resp, "orderBook")), &response); err != nil {
				return nil, err
			} else {
				return response, nil
			}
		} else {
			return nil, fmt.Errorf("failed to get orders information. %v", utils.GetJsonValueAtPath(resp, "message"))
		}
	}
}

func (c *client) GetOrder(orderId string) (api.OrderBook, error) {
	if resp, err := c.invoke(utils.GET, OrdersUrl+fmt.Sprintf(IdQueryParam, orderId), nil); err != nil {
		return api.OrderBook{}, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response api.OrderBook
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(resp, "orderBook")), &response); err != nil {
				return api.OrderBook{}, err
			} else {
				return response, nil
			}
		} else {
			return api.OrderBook{}, fmt.Errorf("failed to get orders information. %v", utils.GetJsonValueAtPath(resp, "message"))
		}
	}
}
