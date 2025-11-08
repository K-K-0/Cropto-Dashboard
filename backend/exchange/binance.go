package exchange

import (
	"bytes"
	"cropto-dashboard/types"

	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func cleanPrice(price types.FlexString) string {
	priceStr := string(price)

	if val, err := strconv.ParseFloat(priceStr, 64); err == nil {
		return fmt.Sprintf("%.8f", val)
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
	log.Println("Connected to Binance WebSocket")

	return nil
}

func (b *BinanceClient) BuildStreamName() string {
	if len(b.Symbols) == 0 {

		return "btcusdt@trade/btcusdt@ticker"
	}

	streams := ""
	for i, symbol := range b.Symbols {
		if i > 0 {
			streams += "/"
		}

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
			log.Printf(" Connection failed: %v. Retrying in %v", err, b.ReconnectDelay)
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
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			if normalized == nil {
				continue
			}

			if bytes.Contains(normalized, []byte(`"eventType":"trade"`)) {
				tradeCount++
			} else if bytes.Contains(normalized, []byte(`"eventType":"ticker"`)) {
				tickerCount++
			}

			if msgCount%50 == 0 {
				log.Printf("Stats - Total: %d | Trades: %d | Tickers: %d", msgCount, tradeCount, tickerCount)
			}

			if debugCount < 10 {
				log.Printf("Sending message #%d: %s", msgCount, string(normalized))
				debugCount++
			}

			select {
			case b.MessageChan <- normalized:
			default:
				log.Println(" Message channel full, dropping message")
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

	var wrapper struct {
		Stream string          `json:"stream"`
		Data   json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(data, &wrapper); err == nil && wrapper.Stream != "" {

		data = wrapper.Data
	}

	var eventCheck map[string]interface{}

	if err := json.Unmarshal(data, &eventCheck); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	eventType := ""
	if e, ok := eventCheck["e"]; ok {
		switch v := e.(type) {
		case string:
			eventType = v
		case float64:
			eventType = fmt.Sprintf("%.0f", v)
		default:

			return nil, nil
		}
	} else {

		return nil, nil
	}

	// Handle TRADE events (real-time price updates)
	if eventType == "trade" {
		var tradeData types.BinanceTradeData

		if err := json.Unmarshal(data, &tradeData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal trade: %w", err)
		}

		ticker := types.TickerMessage{
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

		jsonBytes, _ := json.Marshal(ticker)
		log.Printf("TRADE message: %s", string(jsonBytes))

		return json.Marshal(ticker)
	}

	// Handle 24hrTicker events (statistics)
	if eventType == "24hrTicker" {
		var binanceData types.BinanceTickerData

		if err := json.Unmarshal(data, &binanceData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ticker: %w", err)
		}

		ticker := types.TickerMessage{
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

		jsonBytes, _ := json.Marshal(ticker)
		log.Printf("TICKER message: %s", string(jsonBytes))

		return json.Marshal(ticker)
	}

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
