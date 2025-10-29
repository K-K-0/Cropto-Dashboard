package exchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

// BinanceTickerData represents the ticker message from Binance
// IMPORTANT: JSON keys must match Binance's actual field names (e, E, s, p, P, c, v, q, h, l)
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

type TickerMessage struct {
	Symbol        string `json:"symbol"`
	Price         string `json:"price"`
	Change        string `json:"change"`
	ChangePercent string `json:"changePercent"`
	Volume        string `json:"volume"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Timestamp     int64  `json:"timestamp"`
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
		MessageChan:     make(chan []byte, 256), // Buffered channel
		ReconnectDelay:  1 * time.Second,
		ShouldReconnect: true,
	}
}

const binanceWSURL = "wss://stream.binance.com:9443/stream?streams="

func (c *BinanceClient) Connect() error {
	streamName := c.BuildStreamName()
	url := fmt.Sprintf("%s%s", binanceWSURL, streamName)

	log.Println("Connecting to Binance:", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.Conn = conn
	c.ReconnectDelay = 1 * time.Second
	log.Println("✅ Connected to Binance WebSocket")

	return nil
}

func (b *BinanceClient) BuildStreamName() string {
	if len(b.Symbols) == 0 {
		return "btcusdt@ticker"
	}

	streams := ""
	for i, symbol := range b.Symbols {
		if i > 0 {
			streams += "/"
		}
		streams += fmt.Sprintf("%s@ticker", symbol)
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

			// Exponential backoff
			b.ReconnectDelay *= 2
			if b.ReconnectDelay > 120*time.Second {
				b.ReconnectDelay = 120 * time.Second
			}
			continue
		}

		// Connected successfully, start reading
		b.readLoop()

		// Connection lost, will retry
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

	// Read messages in separate goroutine
	go func() {
		defer close(done)
		for {
			_, message, err := b.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			// Parse and normalize
			normalized, err := b.normalizeMessage(message)
			if err != nil {
				log.Printf("⚠️ Failed to parse message: %v", err)
				continue
			}

			if normalized == nil {
				continue
			}

			// Send to channel (non-blocking)
			select {
			case b.MessageChan <- normalized:
			default:
				log.Println("⚠️  Message channel full, dropping message")
			}
		}
	}()

	// Ping loop
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
	// Trim any whitespace or newlines
	data = bytes.TrimSpace(data)

	// First check if this is a combined stream wrapper
	var wrapper struct {
		Stream string          `json:"stream"`
		Data   json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(data, &wrapper); err == nil && wrapper.Stream != "" {
		// This is wrapped format, extract the data
		data = wrapper.Data
	}

	// Now parse the actual ticker data
	var binanceData BinanceTickerData

	if err := json.Unmarshal(data, &binanceData); err != nil {
		log.Printf("❌ Unmarshal error. Data: %s", string(data[:min(len(data), 200)]))
		return nil, fmt.Errorf("failed to unmarshal ticker: %w", err)
	}

	// Debug: Log what we parsed (first few times)
	log.Printf("✅ Parsed ticker - Event: %s, Symbol: %s, Price: %s",
		binanceData.EventType, binanceData.Symbol, binanceData.LastPrice.String())

	// Skip non-ticker events (this is normal, not an error)
	if binanceData.EventType != "24hrTicker" {
		return nil, nil // Return nil without error
	}

	// Convert to our normalized format
	ticker := TickerMessage{
		Symbol:        binanceData.Symbol,
		Price:         binanceData.LastPrice.String(),
		Change:        binanceData.PriceChange.String(),
		ChangePercent: binanceData.PriceChangePercent.String(),
		Volume:        binanceData.Volume.String(),
		High:          binanceData.HighPrice.String(),
		Low:           binanceData.LowPrice.String(),
		Timestamp:     binanceData.EventTime,
	}

	return json.Marshal(ticker)
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
