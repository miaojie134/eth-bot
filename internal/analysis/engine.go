// internal/analysis/engine.go

package analysis

import (
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/utils"
)

type AnalysisEngine struct {
	indicators []Indicator
}

func NewAnalysisEngine() *AnalysisEngine {
	return &AnalysisEngine{
		indicators: make([]Indicator, 0),
	}
}

func (e *AnalysisEngine) AddIndicator(indicator Indicator) {
	e.indicators = append(e.indicators, indicator)
}

func (e *AnalysisEngine) Analyze(bars []models.Bar) (*AnalysisResult, error) {
	utils.Log.Infof("数据点数量：%d", len(bars))
	result := &AnalysisResult{
		Timestamp:   time.Now(),
		MarketState: Neutral,
		Indicators:  make(map[string]interface{}),
	}

	for _, indicator := range e.indicators {
		value, err := indicator.Calculate(bars)
		if err != nil {
			return nil, err
		}
		result.Indicators[indicator.Name()] = value
	}

	// 简单的市场状态判断逻辑
	if len(bars) > 1 {
		if bars[len(bars)-1].Close > bars[len(bars)-2].Close {
			result.MarketState = Bullish
		} else if bars[len(bars)-1].Close < bars[len(bars)-2].Close {
			result.MarketState = Bearish
		}
	}
	utils.Log.Infof("分析结果：%+v", result)
	return result, nil
}
