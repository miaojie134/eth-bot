// internal/analysis/engine.go
package analysis

import (
	"time"

	"github.com/qqqq/eth-trading-system/internal/analysis/support_resistance"
	"github.com/qqqq/eth-trading-system/internal/analysis/trend"
	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/utils"
)

type AnalysisEngine struct {
	indicators                []Indicator
	trendAnalyzer             *trend.TrendAnalyzer
	supportResistanceAnalyzer *support_resistance.SupportResistanceAnalyzer
}

func NewAnalysisEngine() Engine {
	return &AnalysisEngine{
		indicators:                make([]Indicator, 0),
		trendAnalyzer:             trend.NewTrendAnalyzer(10, 30),
		supportResistanceAnalyzer: support_resistance.NewSupportResistanceAnalyzer(20),
	}
}

func (e *AnalysisEngine) AddIndicator(indicator Indicator) {
	e.indicators = append(e.indicators, indicator)
}

func (e *AnalysisEngine) Analyze(bars []models.Bar) (*models.AnalysisResult, error) {
	utils.Log.Infof("数据点数量：%d", len(bars))
	result := &models.AnalysisResult{
		Timestamp:   time.Now(),
		MarketState: models.Neutral,
		Indicators:  make(map[string]interface{}),
	}

	for _, indicator := range e.indicators {
		value, err := indicator.Calculate(bars)
		if err != nil {
			return nil, err
		}
		result.Indicators[indicator.Name()] = value
	}

	trendState, err := e.trendAnalyzer.AnalyzeTrend(bars)
	if err != nil {
		return nil, err
	}
	result.Indicators["Trend"] = trendState

	supports, resistances := e.supportResistanceAnalyzer.FindLevels(bars)
	result.Indicators["Supports"] = supports
	result.Indicators["Resistances"] = resistances

	// 市场状态判断逻辑
	if trendState == "Uptrend" {
		result.MarketState = models.Bullish
	} else if trendState == "Downtrend" {
		result.MarketState = models.Bearish
	}

	utils.Log.Infof("分析结果：%+v", result)
	return result, nil
}
