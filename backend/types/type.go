package types

type BinanceTickerData struct {
	EventType          string     `json:"e"`
	EventTime          int64      `json:"E"`
	Symbol             string     `json:"s"`
	PriceChange        FlexString `json:"p"`
	PriceChangePercent FlexString `json:"P"`
	LastPrice          FlexString `json:"c"`
	Volume             FlexString `json:"v"`
	QuoteVolume        FlexString `json:"q"`
	HighPrice          FlexString `json:"h"`
	LowPrice           FlexString `json:"l"`
}

type BinanceTradeData struct {
	EventType     string     `json:"e"`
	EventTime     int64      `json:"E"`
	Symbol        string     `json:"s"`
	TradeID       int64      `json:"t"`
	Price         FlexString `json:"p"`
	Quantity      FlexString `json:"q"`
	BuyerOrderID  int64      `json:"b"`
	SellerOrderID int64      `json:"a"`
	TradeTime     int64      `json:"T"`
	IsBuyerMaker  bool       `json:"m"`
}

type TickerMessage struct {
	Symbol        string `json:"symbol"`
	Price         string `json:"price"`
	Change        string `json:"change"`
	ChangePercent string `json:"changePercent"`
	Volume        string `json:"volume"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Timestamp     int64  `json:"timestamp"`
	EventType     string `json:"eventType"`
}
