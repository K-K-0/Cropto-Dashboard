package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const baseUrl = "https://api.binance.com/api/v3"

var httpClient = &http.Client{Timeout: 10 * time.Second}

type CandleStick struct {
	OpenTime  int64  `json:"openTime"`
	Open      string `json:"open"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Close     string `json:"close"`
	Volume    string `json:"volume"`
	CloseTime int64  `json:"closeTime"`
}

type TechnicalIndicators struct {
	MA20      []float64 `json:"ma20,omitempty"`
	MA50      []float64 `json:"ma50,omitempty"`
	MA200     []float64 `json:"ma200,omitempty"`
	RSI       []float64 `json:"rsi,omitempty"`
	MACD      []float64 `json:"macd,omitempty"`
	Signal    []float64 `json:"signal,omitempty"`
	Histogram []float64 `json:"histogram,omitempty"`
}

type ChartData struct {
	Symbol       string              `json:"symbol"`
	Interval     string              `json:"interval"`
	Candlesticks []CandleStick       `json:"candlesticks"`
	Indicators   TechnicalIndicators `json:"indicators"`
}

func fetch(url string) ([]byte, error) {
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("binance error: %s", body)
	}

	return io.ReadAll(res.Body)
}

func GetHistoricalData(symbol, interval string, limit int) (*ChartData, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol)) // FIX 1

	if limit <= 0 || limit > 1000 {
		limit = 500
	}

	url := fmt.Sprintf("%s/klines?symbol=%s&interval=%s&limit=%d",
		baseUrl, symbol, interval, limit)

	body, err := fetch(url)
	if err != nil {
		return nil, err
	}

	var raw [][]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	candles := make([]CandleStick, 0, len(raw))

	for _, r := range raw {
		candles = append(candles, CandleStick{
			OpenTime:  int64(r[0].(float64)),
			Open:      fmt.Sprint(r[1]),
			High:      fmt.Sprint(r[2]),
			Low:       fmt.Sprint(r[3]),
			Close:     fmt.Sprint(r[4]),
			Volume:    fmt.Sprint(r[5]),
			CloseTime: int64(r[6].(float64)),
		})
	}

	return &ChartData{
		Symbol:       symbol,
		Interval:     interval,
		Candlesticks: candles,
	}, nil
}

func GetMultipleHistoricalData(symbols []string, interval string, limit int) (map[string]*ChartData, error) {
	result := make(map[string]*ChartData)

	for _, s := range symbols {
		data, err := GetHistoricalData(s, interval, limit)
		if err == nil {
			result[s] = data
		}
		time.Sleep(50 * time.Millisecond)
	}

	return result, nil
}

func GetLatestPrice(symbol string) (string, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	url := fmt.Sprintf("%s/ticker/price?symbol=%s", baseUrl, symbol)
	body, err := fetch(url)
	if err != nil {
		return "", err
	}

	var res struct {
		Price string `json:"price"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res.Price, nil
}
