package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rishi-anand/fyers-go-client/utils"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	apiKey      string
	accessToken string
	debug       bool
	httpClient  *http.Client
}

func New(apiKey, accessToken string) *Client {
	return &Client{
		apiKey:      apiKey,
		accessToken: accessToken,
		debug:       false,
		httpClient: &http.Client{
			Timeout: time.Duration(20) * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   10,
				ResponseHeaderTimeout: time.Second * time.Duration(20),
			},
		},
	}
}

func (c *Client) WithHttpClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

func (c *Client) EnableDebug() *Client {
	c.debug = true
	return c
}

func (c *Client) invoke(method, url string, body interface{}) ([]byte, error) {
	headerMap := map[string]string{
		"Authorization": fmt.Sprintf("%s:%s", c.apiKey, c.accessToken),
		"Content-Type":  "application/json",
	}

	var bodyByte []byte
	if bodyByteArr, err := json.Marshal(body); err != nil {
		return nil, err
	} else {
		bodyByte = bodyByteArr
	}

	if resp, err := utils.DoHttpCall(method, url, bodyByte, headerMap); err != nil {
		log.Error("Failed to make http call", err)
		return nil, err
	} else {
		return resp, nil
	}
}
