package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

func GetHistoricalDataWithIndicators(symbol, interval string, limit int, includeIndicators bool) (*ChartData, error) {
	data, err := GetHistoricalData(symbol, interval, limit)
	if err != nil {
		return nil, err
	}

	if !includeIndicators || len(data.Candlesticks) == 0 {
		return data, nil
	}

	closePrice := make([]float64, len(data.Candlesticks))
	highPrice := make([]float64, len(data.Candlesticks))
	lowPrice := make([]float64, len(data.Candlesticks))

	for i, candle := range data.Candlesticks {
		close, _ := strconv.ParseFloat(candle.Close, 64)
		high, _ := strconv.ParseFloat(candle.High, 64)
		low, _ := strconv.ParseFloat(candle.Low, 64)
		closePrice[i] = close
		highPrice[i] = high
		lowPrice[i] = low
	}

	indicators := TechnicalIndicators{}

	if len(closePrice) >= 20 {
		indicators.MA20 = calculateMA(closePrice, 20)
	}

	if len(closePrice) >= 50 {
		indicators.MA50 = calculateMA(closePrice, 50)
	}

	if len(closePrice) >= 200 {
		indicators.MA200 = calculateMA(closePrice, 200)
	}

	if len(closePrice) >= 14 {
		indicators.RSI = calculateMA(closePrice, 14)
	}

	if len(closePrice) >= 26 {
		macd, signal, histogram := calculateMACD(closePrice, 12, 26, 9)
		indicators.MACD = macd
		indicators.Signal = signal
		indicators.Histogram = histogram
	}
	data.Indicators = indicators
	return data, nil
}

func calculateMA(prices []float64, period int) []float64 {
	ma := make([]float64, len(prices))
	for i := period - 1; i < len(prices); i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i-j]
		}
		ma[i] = sum / float64(period)
	}

	return ma
}

func calculateRSI(prices []float64, period int) []float64 {
	rsi := make([]float64, len(prices))
	gains := make([]float64, len(prices))
	losses := make([]float64, len(prices))

	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = -change
		}
	}

	avgGain := 0.0
	avgLoss := 0.0
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	for i := period; i < len(prices); i++ {
		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)
	}
	return rsi
}

func calculateEMA(prices []float64, period int) []float64 {
	ema := make([]float64, len(prices))
	multiplier := 2.0 / float64(period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i]-ema[i-1])*multiplier + ema[i-1]
	}
	return ema
}

func calculateMACD(price []float64, fast, slow, signal int) ([]float64, []float64, []float64) {
	fastEMA := calculateEMA(price, fast)
	slowEMA := calculateEMA(price, slow)

	macd := make([]float64, len(price))

	for i := 0; i < len(price); i++ {
		macd[i] = fastEMA[i] - slowEMA[i]
	}

	signalLine := calculateEMA(macd, signal)
	histogram := make([]float64, len(price))
	for i := 0; i < len(price); i++ {
		histogram[i] = macd[i] - signalLine[i]
	}

	return macd, signalLine, histogram
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
