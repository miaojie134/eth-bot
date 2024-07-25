package strategy

import (
	"github.com/qqqq/eth-trading-system/internal/analysis"
	"github.com/qqqq/eth-trading-system/internal/models"
)

type SimpleMAStrategy struct {
	BaseStrategy
	ShortPeriod int
	LongPeriod  int
}

func NewSimpleMAStrategy(shortPeriod, longPeriod int) *SimpleMAStrategy {
	return &SimpleMAStrategy{
		BaseStrategy: BaseStrategy{name: "SimpleMA"},
		ShortPeriod:  shortPeriod,
		LongPeriod:   longPeriod,
	}
}

func (s *SimpleMAStrategy) Evaluate(data []models.Bar, analysisResult *models.AnalysisResult) *models.TradeSignal {
	shortMA, ok := analysisResult.Indicators[analysis.IndicatorSMA(s.ShortPeriod)].([]float64)
	if !ok {
		return &models.TradeSignal{StrategyName: s.Name(), Action: "HOLD", Reason: "短期MA指标不可用"}
	}

	longMA, ok := analysisResult.Indicators[analysis.IndicatorSMA(s.LongPeriod)].([]float64)
	if !ok {
		return &models.TradeSignal{StrategyName: s.Name(), Action: "HOLD", Reason: "长期MA指标不可用"}
	}

	if len(shortMA) < 2 || len(longMA) < 2 {
		return &models.TradeSignal{StrategyName: s.Name(), Action: "HOLD", Reason: "不足够的数据"}
	}

	if shortMA[len(shortMA)-1] > longMA[len(longMA)-1] &&
		shortMA[len(shortMA)-2] <= longMA[len(longMA)-2] {
		return &models.TradeSignal{
			StrategyName: s.Name(),
			Action:       "BUY",
			Price:        data[len(data)-1].Close,
			Reason:       "短期均线上穿长期均线",
		}
	} else if shortMA[len(shortMA)-1] < longMA[len(longMA)-1] &&
		shortMA[len(shortMA)-2] >= longMA[len(longMA)-2] {
		return &models.TradeSignal{
			StrategyName: s.Name(),
			Action:       "SELL",
			Price:        data[len(data)-1].Close,
			Reason:       "短期均线下穿长期均线",
		}
	}

	return &models.TradeSignal{StrategyName: s.Name(), Action: "HOLD", Reason: "无显著变化"}
}
