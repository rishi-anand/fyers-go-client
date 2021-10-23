package fyers

const (
	Host             = "https://api.fyers.in"
	SingleOrderUrl   = Host + "/api/v2/orders"
	MultiOrderUrl    = Host + "/api/v2/orders-multi"
	PositionOrderUrl = Host + "/api/v2/positions"
	QuoteUrl         = Host + "/data-rest/v2/quotes/?symbols=%s"
)
