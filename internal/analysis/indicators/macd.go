package indicators

import (
	"fmt"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type MACD struct {
	FastPeriod   int
	SlowPeriod   int
	SignalPeriod int
}

func NewMACD(fastPeriod, slowPeriod, signalPeriod int) *MACD {
	return &MACD{
		FastPeriod:   fastPeriod,
		SlowPeriod:   slowPeriod,
		SignalPeriod: signalPeriod,
	}
}

func (macd *MACD) Calculate(bars []models.Bar) (interface{}, error) {
	if len(bars) < macd.SlowPeriod {
		return nil, fmt.Errorf("not enough data for MACD calculation")
	}

	fastEMA := calculateEMA(bars, macd.FastPeriod)
	slowEMA := calculateEMA(bars, macd.SlowPeriod)

	macdLine := fastEMA[len(fastEMA)-1] - slowEMA[len(slowEMA)-1]
	signalLine := calculateEMA(bars[len(bars)-len(slowEMA):], macd.SignalPeriod)[macd.SignalPeriod-1]
	histogram := macdLine - signalLine

	return map[string]float64{
		"MACD":      macdLine,
		"Signal":    signalLine,
		"Histogram": histogram,
	}, nil
}

func (macd *MACD) Name() string {
	return fmt.Sprintf("MACD(%d,%d,%d)", macd.FastPeriod, macd.SlowPeriod, macd.SignalPeriod)
}

func calculateEMA(bars []models.Bar, period int) []float64 {
	ema := make([]float64, len(bars))
	k := 2.0 / float64(period+1)

	for i := range bars {
		if i == 0 {
			ema[i] = bars[i].Close
		} else {
			ema[i] = bars[i].Close*k + ema[i-1]*(1-k)
		}
	}

	return ema
}
