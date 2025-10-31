package exchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// FlexString handles both string and number types from JSON
type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexString(s)
		return nil
	}

	// If that fails, try as number
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexString(fmt.Sprintf("%f", num))
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexString", data)
}

func (f FlexString) String() string {
	return string(f)
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BinanceTickerData represents the 24hr ticker message from Binance
type BinanceTickerData struct {
	EventType          string     `json:"e"` // "24hrTicker"
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

// BinanceTradeData represents individual trade messages (real-time)
type BinanceTradeData struct {
	EventType     string     `json:"e"` // "trade"
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
	EventType     string `json:"eventType"` // "trade" or "ticker"
}

// Helper to clean price strings (remove trailing zeros from float conversion)
func cleanPrice(price FlexString) string {
	priceStr := price.String()
	// Remove excessive decimal places from float conversion
	// "67234.500000" -> "67234.50"
	if val, err := strconv.ParseFloat(priceStr, 64); err == nil {
		return fmt.Sprintf("%.8f", val) // Keep 8 decimals, trim trailing zeros later
	}
	return priceStr
}

type BinanceClient struct {
	Symbols         []string
	Conn            *websocket.Conn
	MessageChan     chan []byte
	ReconnectDelay  time.Duration
	ShouldReconnect bool
}

func NewBinanceClient(symbols []string) *BinanceClient {
	return &BinanceClient{
		Symbols:         symbols,
		MessageChan:     make(chan []byte, 256),
		ReconnectDelay:  1 * time.Second,
		ShouldReconnect: true,
	}
}

const binanceWSURL = "wss://stream.binance.com:9443/stream?streams="

func (c *BinanceClient) Connect() error {
	streamName := c.BuildStreamName()
	url := fmt.Sprintf("%s%s", binanceWSURL, streamName)

	log.Println("🔗 Connecting to Binance:", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.Conn = conn
	c.ReconnectDelay = 1 * time.Second
	log.Println("✅ Connected to Binance WebSocket")

	return nil
}

// BuildStreamName creates subscription for BOTH trade and ticker streams
func (b *BinanceClient) BuildStreamName() string {
	if len(b.Symbols) == 0 {
		// Default: both trade and ticker for BTC
		return "btcusdt@trade/btcusdt@ticker"
	}

	streams := ""
	for i, symbol := range b.Symbols {
		if i > 0 {
			streams += "/"
		}
		// Subscribe to BOTH trade (live price) and ticker (24h stats)
		streams += fmt.Sprintf("%s@trade/%s@ticker", symbol, symbol)
	}
	return streams
}

func (b *BinanceClient) Start() {
	go b.ReconnectLoop()
}

func (b *BinanceClient) ReconnectLoop() {
	for b.ShouldReconnect {
		err := b.Connect()
		if err != nil {
			log.Printf("❌ Connection failed: %v. Retrying in %v", err, b.ReconnectDelay)
			time.Sleep(b.ReconnectDelay)

			b.ReconnectDelay *= 2
			if b.ReconnectDelay > 120*time.Second {
				b.ReconnectDelay = 120 * time.Second
			}
			continue
		}

		b.readLoop()

		log.Println("🔌 Connection lost, reconnecting...")
		time.Sleep(b.ReconnectDelay)
	}
}

func (b *BinanceClient) readLoop() {
	defer func() {
		if b.Conn != nil {
			b.Conn.Close()
		}
	}()

	b.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	b.Conn.SetPongHandler(func(string) error {
		b.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})

	go func() {
		defer close(done)
		msgCount := 0
		debugCount := 0
		tradeCount := 0
		tickerCount := 0

		for {
			_, message, err := b.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			msgCount++

			normalized, err := b.normalizeMessage(message)
			if err != nil {
				log.Printf("⚠️ Failed to parse message: %v", err)
				continue
			}

			if normalized == nil {
				continue
			}

			// Count event types
			if bytes.Contains(normalized, []byte(`"eventType":"trade"`)) {
				tradeCount++
			} else if bytes.Contains(normalized, []byte(`"eventType":"ticker"`)) {
				tickerCount++
			}

			// Log stats every 50 messages
			if msgCount%50 == 0 {
				log.Printf("📊 Stats - Total: %d | Trades: %d | Tickers: %d", msgCount, tradeCount, tickerCount)
			}

			// Only log first 10 messages
			if debugCount < 10 {
				log.Printf("✅ Sending message #%d: %s", msgCount, string(normalized))
				debugCount++
			}

			select {
			case b.MessageChan <- normalized:
			default:
				log.Println("⚠️  Message channel full, dropping message")
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := b.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Printf("Ping failed: %v", err)
				return
			}
		}
	}
}

func (b *BinanceClient) normalizeMessage(data []byte) ([]byte, error) {
	data = bytes.TrimSpace(data)

	// Check if this is a combined stream wrapper
	var wrapper struct {
		Stream string          `json:"stream"`
		Data   json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(data, &wrapper); err == nil && wrapper.Stream != "" {
		// Extract the actual data from wrapper
		data = wrapper.Data
	}

	// Determine event type - use flexible map to handle string OR number
	var eventCheck map[string]interface{}

	if err := json.Unmarshal(data, &eventCheck); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Get event type (handle both string and number)
	eventType := ""
	if e, ok := eventCheck["e"]; ok {
		switch v := e.(type) {
		case string:
			eventType = v
		case float64:
			eventType = fmt.Sprintf("%.0f", v)
		default:
			// Unknown type, skip
			return nil, nil
		}
	} else {
		// No event type field, skip
		return nil, nil
	}

	// Handle TRADE events (real-time price updates)
	if eventType == "trade" {
		var tradeData BinanceTradeData

		if err := json.Unmarshal(data, &tradeData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal trade: %w", err)
		}

		// Convert to our standard format - EXPLICITLY set each field
		ticker := TickerMessage{
			Symbol:        tradeData.Symbol,
			Price:         cleanPrice(tradeData.Price),
			Change:        "0",
			ChangePercent: "0",
			Volume:        cleanPrice(tradeData.Quantity),
			High:          "0",
			Low:           "0",
			Timestamp:     tradeData.TradeTime,
			EventType:     "trade",
		}

		// Debug first few
		jsonBytes, _ := json.Marshal(ticker)
		log.Printf("🔵 TRADE message: %s", string(jsonBytes))

		return json.Marshal(ticker)
	}

	// Handle 24hrTicker events (statistics)
	if eventType == "24hrTicker" {
		var binanceData BinanceTickerData

		if err := json.Unmarshal(data, &binanceData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ticker: %w", err)
		}

		// Convert to our standard format - EXPLICITLY set each field
		ticker := TickerMessage{
			Symbol:        binanceData.Symbol,
			Price:         cleanPrice(binanceData.LastPrice),
			Change:        cleanPrice(binanceData.PriceChange),
			ChangePercent: cleanPrice(binanceData.PriceChangePercent),
			Volume:        cleanPrice(binanceData.Volume),
			High:          cleanPrice(binanceData.HighPrice),
			Low:           cleanPrice(binanceData.LowPrice),
			Timestamp:     binanceData.EventTime,
			EventType:     "ticker",
		}

		// Debug first few
		jsonBytes, _ := json.Marshal(ticker)
		log.Printf("🟢 TICKER message: %s", string(jsonBytes))

		return json.Marshal(ticker)
	}

	// Ignore other event types silently
	return nil, nil
}

func (b *BinanceClient) GetMessageChannel() <-chan []byte {
	return b.MessageChan
}

func (b *BinanceClient) Close() {
	b.ShouldReconnect = false
	if b.Conn != nil {
		b.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		b.Conn.Close()
	}
	close(b.MessageChan)
}
