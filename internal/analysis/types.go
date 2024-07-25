// internal/analysis/types.go

package analysis

import (
	"fmt"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type Indicator interface {
	Calculate(bars []models.Bar) (interface{}, error)
	Name() string
}

type Engine interface {
	Analyze(bars []models.Bar) (*models.AnalysisResult, error)
	AddIndicator(indicator Indicator)
}

const (
	IndicatorSMAPrefix     = "SMA"
	IndicatorMACD          = "MACD"
	IndicatorMACDSignal    = "MACDSignal"
	IndicatorMACDHistogram = "MACDHistogram"
)

// IndicatorSMA returns the name of the Simple Moving Average indicator for a given period
func IndicatorSMA(period int) string {
	return fmt.Sprintf("%s%d", IndicatorSMAPrefix, period)
}
