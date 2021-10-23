# Fyers API GoLang Client

[![Go Reference](https://pkg.go.dev/badge/github.com/rishi-anand/fyers-go-client.svg)](https://pkg.go.dev/github.com/rishi-anand/fyers-go-client)

Fyers api golang client will help user to connect with various apis and subscribing with websocket notification. User can execute orders at real time and get order status as well as stream live market data via websocket.

## Features

- Place real time single and multiple orders
- Get quotes of single of multiple symbols
- Stream live data from market via fyers go client
- More to be added soon.. 😊

## Installation

Fyers client is just a `go get` away
```sh
go get github.com/rishi-anand/fyers-go-client
```

## API usage

```
package main

import (
	"fmt"

	"github.com/rishi-anand/fyers-go-client"
)

func main() {
	apiKey := "<YOUR_API_KEY>"
	accessToken := "<YOUR_ACCESS_TOKEN>"
	symbols := []string{"NSE:SBIN-EQ", "NSE:ONGC-EQ"}

	cli := fyers.New(apiKey, accessToken)
	if quotes, err := cli.GetQuote(symbols); err != nil {
		fmt.Errorf("failed to get quote from fyers. %v", err)
	} else {
		fmt.Println(quotes)
	}
}
```

## Websocket Notification (for live market data)

```
package main

import (
	"fmt"

	"github.com/rishi-anand/fyers-go-client/api"
	fyerswatch "github.com/rishi-anand/fyers-go-client/websocket"
)

func main() {
	apiKey := "<YOUR_API_KEY>"
	accessToken := "<YOUR_ACCESS_TOKEN>"

	onConnectFunc := func() {
		fmt.Println("watch subscription is connected")
	}

	onMessageFunc := func(notification api.Notification) {
		fmt.Println(notification.Type, notification.SymbolData)
	}

	onErrorFunc := func(err error) {
		fmt.Errorf("failed to watch | disconnected from watch. %v", err)
	}

	onCloseFunc := func() {
		fmt.Println("watch connection is closed")
	}

	cli := fyerswatch.NewNotifier(apiKey, accessToken).
		WithOnConnectFunc(onConnectFunc).
		WithOnMessageFunc(onMessageFunc).
		WithOnErrorFunc(onErrorFunc).
		WithOnCloseFunc(onCloseFunc)

	cli.Subscribe(api.SymbolDataTick, "NSE:SBIN-EQ", "NSE:ONGC-EQ")

	/*
		symbols := []string {"NSE:SBIN-EQ", "NSE:ONGC-EQ"}
		cli.Subscribe(api.SymbolDataTick, symbols...)
	*/
}
```

## Examples

Check examples [examples](https://github.com/rishi-anand/fyers-go-client/tree/main/examples) folder for more examples.

