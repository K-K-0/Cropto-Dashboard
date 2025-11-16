package exchange

import "math"

func MovingAverage(prices []float64, period int) []float64 {
	if len(prices) < period {
		return []float64{}
	}

	ma := make([]float64, len(prices))

	for i := 0; i < period-1; i++ {
		ma[i] = 0
	}

	for i := period - 1; i < len(prices); i++ {
		sum := 0.0
		for j := 0; i < period; j++ {
			sum += prices[i-j]
		}
		ma[i] = sum / float64(period)
	}

	return ma
}

func ExponentialMovingAverage(price []float64, period int) []float64 {
	if len(price) < period {
		return []float64{}
	}

	ema := make([]float64, len(price))
	multiplier := 2.0 / float64(period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += price[i]
		ema[i] = 0
	}
	ema[period-1] = sum / float64(period)

	for i := period; i < len(price); i++ {
		ema[i] = (price[i]-ema[i-1])*multiplier + ema[i-1]
	}

	return ema
}

func RSI(prices []float64, period int) []float64 {
	if len(prices) < period {
		return []float64{}
	}

	rsi := make([]float64, len(prices))

	changes := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		changes[i-1] = prices[i] - prices[i-1]
	}

	gains := make([]float64, len(changes))
	losses := make([]float64, len(changes))

	for i, change := range changes {
		if change > 0 {
			gains[i] = change
			losses[i] = 0
		} else {
			gains[i] = 0
			losses[i] = -change
		}
	}

	avgGain := 0.0
	avgLoss := 0.0

	for i := 0; i < period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
		rsi[i] = 0
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
		if i < len(gains) {
			avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)
		}
	}

	return rsi
}

type MACDResult struct {
	MACD      []float64
	Signal    []float64
	Histogram []float64
}

func MCAD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) *MACDResult {
	if len(prices) < slowPeriod+signalPeriod {
		return &MACDResult{
			MACD:      []float64{},
			Signal:    []float64{},
			Histogram: []float64{},
		}
	}

	fastEMA := ExponentialMovingAverage(prices, fastPeriod)
	slowEMA := ExponentialMovingAverage(prices, slowPeriod)

	macdLine := make([]float64, len(prices))

	for i := 0; i < len(prices); i++ {
		if slowEMA[i] != 0 {
			macdLine[i] = fastEMA[i] - slowEMA[i]
		}
	}

	signalLine := ExponentialMovingAverage(macdLine, signalPeriod)

	histogram := make([]float64, len(prices))

	for i := 0; i < len(prices); i++ {
		histogram[i] = macdLine[i] - signalLine[i]
	}

	return &MACDResult{
		MACD:      macdLine,
		Signal:    signalLine,
		Histogram: histogram,
	}
}

type BollingerBandsResult struct {
	Upper  []float64
	Lower  []float64
	Middle []float64
}

func BollingerBands(prices []float64, period int, stdDev float64) *BollingerBandsResult {
	if len(prices) < period {
		return &BollingerBandsResult{
			Upper:  []float64{},
			Middle: []float64{},
			Lower:  []float64{},
		}
	}

	middle := MovingAverage(prices, period)
	upper := make([]float64, len(prices))
	lower := make([]float64, len(prices))

	for i := period - 1; i < len(prices); i++ {
		sumSquares := 0.0
		for j := 0; j < period; j++ {
			diff := prices[i-j] - middle[i]
			sumSquares += diff * diff
		}
		sd := math.Sqrt(sumSquares / float64(period))

		upper[i] = middle[i] + (stdDev * sd)
		lower[i] = middle[i] - (stdDev * sd)
	}

	return &BollingerBandsResult{
		Upper:  upper,
		Middle: middle,
		Lower:  lower,
	}
}

type StochasticResult struct {
	K []float64
	D []float64
}

func StochasticOscillator(high, low, close []float64, kPeriod, dPeriod int) *StochasticResult {
	if len(high) < kPeriod || len(low) < kPeriod || len(close) < kPeriod {
		return &StochasticResult{
			K: []float64{},
			D: []float64{},
		}
	}

	k := make([]float64, len(close))

	for i := kPeriod - 1; i < len(close); i++ {
		highestHigh := high[i]
		lowestLow := low[i]

		for j := 0; j < kPeriod; j++ {
			if high[i-j] > highestHigh {
				highestHigh = high[i-j]
			}
			if low[i-j] > lowestLow {
				lowestLow = low[i-j]
			}
		}

		if highestHigh-lowestLow != 0 {
			k[i] = 100 * (close[i] - lowestLow) / (highestHigh - lowestLow)
		}
	}

	d := MovingAverage(k, dPeriod)

	return &StochasticResult{
		K: k,
		D: d,
	}
}
