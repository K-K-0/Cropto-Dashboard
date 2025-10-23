package exchange

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const binanceWSURL = "wss://stream.binance.com:9443/ws"

type BinanceTickerData struct {
	EventType          string `json:"event_type"`
	EventTime          int64  `json:"event_time"`
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"price_change"`
	PriceChangePercent string `json:"price_change_percentage"`
	LastPrice          string `json:"last_price"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quote_volume"`
	HighPrice          string `json:"high_price"`
	LowPrice           string `json:"low_price"`
}

type TickerMessage struct {
	Symbol        string `json:"symbol"`
	Price         string `json:"price"`
	Change        string `json:"change"`
	ChangePercent string `json:"change_percent"`
	Volume        string `json:"volume"`
	High          string `json:"high"`
	Low           string `json:"low"`
	TimeStamp     int64  `json:"timestamp"`
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
		MessageChan:     make(chan []byte),
		ReconnectDelay:  1 * time.Second,
		ShouldReconnect: true,
	}
}

func (c *BinanceClient) Connect() error {
	streamName := c.BuildStreamName()
	url := fmt.Sprintf("%s/%s", binanceWSURL, streamName)

	log.Println("connected to Binance: ", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("Failed to connect: ", err)
	}

	c.Conn = conn
	c.ReconnectDelay = 1 * time.Second
	log.Println("connected binance websocket")

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
			log.Printf("Connection failed: %v. Reconnect in %v.", err, b.ReconnectDelay)
			time.Sleep(b.ReconnectDelay)

			b.ReconnectDelay *= 2
			if b.ReconnectDelay > 120*time.Second {
				b.ReconnectDelay = 120 * time.Second
			}
			continue
		}
		b.readLoop()
		log.Println("Connection lost, reconnecting......")
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

		_, msg, err := b.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("websocket error: ", err)
			}
			return
		}
		normalized, err := b.normalizeMessage(msg)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		select {
		case b.MessageChan <- normalized:
		default:
			log.Println("Message channel full, dropping message")
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
	var BinanceData BinanceTickerData

	if err := json.Unmarshal(data, &BinanceData); err != nil {
		return nil, err
	}

	if BinanceData.EventType != "24hrTicker" {
		return nil, fmt.Errorf("not a ticker event")
	}

	ticker := TickerMessage{
		Symbol:        BinanceData.Symbol,
		Price:         BinanceData.LastPrice,
		Change:        BinanceData.PriceChange,
		ChangePercent: BinanceData.PriceChangePercent,
		Volume:        BinanceData.Volume,
		High:          BinanceData.HighPrice,
		Low:           BinanceData.LowPrice,
		TimeStamp:     BinanceData.EventTime,
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
