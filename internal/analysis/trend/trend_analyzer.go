package trend

import (
	"github.com/qqqq/eth-trading-system/internal/analysis/indicators"
	"github.com/qqqq/eth-trading-system/internal/models"
)

type TrendAnalyzer struct {
	ShortSMA *indicators.SimpleMovingAverage
	LongSMA  *indicators.SimpleMovingAverage
}

func NewTrendAnalyzer(shortPeriod, longPeriod int) *TrendAnalyzer {
	return &TrendAnalyzer{
		ShortSMA: indicators.NewSimpleMovingAverage(shortPeriod),
		LongSMA:  indicators.NewSimpleMovingAverage(longPeriod),
	}
}

func (ta *TrendAnalyzer) AnalyzeTrend(bars []models.Bar) (string, error) {
	shortSMA, err := ta.ShortSMA.Calculate(bars)
	if err != nil {
		return "", err
	}

	longSMA, err := ta.LongSMA.Calculate(bars)
	if err != nil {
		return "", err
	}

	if shortSMA.(float64) > longSMA.(float64) {
		return "Uptrend", nil
	} else if shortSMA.(float64) < longSMA.(float64) {
		return "Downtrend", nil
	} else {
		return "Neutral", nil
	}
}
