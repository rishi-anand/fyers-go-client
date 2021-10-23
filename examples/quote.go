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
