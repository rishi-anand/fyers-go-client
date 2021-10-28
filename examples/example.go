package main

import (
	"fmt"
	"time"

	"github.com/rishi-anand/fyers-go-client/api"
	fyerswatch "github.com/rishi-anand/fyers-go-client/websocket"

	"github.com/rishi-anand/fyers-go-client"
)

const (
	apiKey      = "<YOUR_API_KEY>"
	accessToken = "<YOUR_ACCESS_TOKEN>"
)

func main() {
	printUserProfile()
	//printQuote()
	//printHistoricalData()
	//subscribe()

}

func printUserProfile() {
	cli := fyers.New(apiKey, accessToken)
	if profile, err := cli.GetProfile(); err != nil {
		fmt.Errorf("failed to get profile from fyers. %v", err)
	} else {
		fmt.Println(profile)
	}
}

func printQuote() {
	symbols := []string{"NSE:HDFCBANK-EQ", "NSE:HDFC-EQ"}
	cli := fyers.New(apiKey, accessToken)
	if quotes, err := cli.GetQuote(symbols); err != nil {
		fmt.Errorf("failed to get quote from fyers. %v", err)
	} else {
		fmt.Println(quotes)
	}
}

func printHistoricalData() {
	/*
		Start Date and End date are required parameters, if not provided it will be defaulted to [2021-01-01, 2021-01-02] .
		Libraries like https://github.com/araddon/dateparse may help user to parse
		date from their desired format to required ones.
	*/
	startDate, err := time.Parse(time.RFC3339, "2021-01-12T11:45:26.371Z")
	if err != nil {
		fmt.Errorf("failed to parse date %v", err)
		return
	}

	endtDate, err := time.Parse(time.RFC3339, "2021-01-15T11:45:26.371Z")
	if err != nil {
		fmt.Errorf("failed to parse date %v", err)
		return
	}
	cli := fyers.New(apiKey, accessToken)
	if data, err := cli.GetHistoricalData("NSE:HDFCBANK-EQ", api.Minute1, startDate, endtDate); err != nil {
		fmt.Errorf("failed to parse date %v", err)
	} else {
		fmt.Println(data)
	}
}

func subscribe() {
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
