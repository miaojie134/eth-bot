// internal/analysis/types.go

package analysis

import (
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type Indicator interface {
	Calculate(bars []models.Bar) (interface{}, error)
	Name() string
}

type MarketState int

const (
	Bullish MarketState = iota
	Bearish
	Neutral
)

type AnalysisResult struct {
	Timestamp   time.Time
	MarketState MarketState
	Indicators  map[string]interface{}
}

type Engine interface {
	Analyze(bars []models.Bar) (*AnalysisResult, error)
	AddIndicator(indicator Indicator)
}
