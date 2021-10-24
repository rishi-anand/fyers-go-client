package fyers

const (
	Host                 = "https://api.fyers.in"
	ApiV2                = "/api/v2"
	OrdersUrl            = Host + ApiV2 + "/orders"
	MultiOrderUrl        = Host + ApiV2 + "/orders-multi"
	PositionsUrl         = Host + ApiV2 + "/positions"
	QuoteUrl             = Host + "/data-rest/v2/quotes/?symbols=%s"
	ProfileUrl           = Host + ApiV2 + "/profile"
	FundsUrl             = Host + ApiV2 + "/funds"
	HoldingsUrl          = Host + ApiV2 + "/holdings"
	TradeBookUrl         = Host + ApiV2 + "/tradebook"
	HistoricalDataApiUrl = Host + "/data-rest/v2/history"
	MarketDepthApiUrl    = Host + "/data-rest/v2/depth/?symbol=%s&ohlcv_flag=1"

	IdQueryParam = "?id=%s"
)
