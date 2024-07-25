package strategy

import (
	"github.com/qqqq/eth-trading-system/internal/models"
)

// CompositeStrategy 组合多个策略的结果
type CompositeStrategy struct {
	BaseStrategy
	strategies []Strategy
	weights    map[string]float64
}

func NewCompositeStrategy(strategies []Strategy, weights map[string]float64) *CompositeStrategy {
	return &CompositeStrategy{
		BaseStrategy: BaseStrategy{name: "Composite"},
		strategies:   strategies,
		weights:      weights,
	}
}

func (cs *CompositeStrategy) Evaluate(data []models.Bar, analysisResult *models.AnalysisResult) *models.TradeSignal {
	var buyScore, sellScore float64
	totalWeight := 0.0

	for _, strategy := range cs.strategies {
		signal := strategy.Evaluate(data, analysisResult)
		weight := cs.weights[strategy.Name()]
		totalWeight += weight

		switch signal.Action {
		case "BUY":
			buyScore += weight
		case "SELL":
			sellScore += weight
		}
	}

	// 归一化得分
	buyScore /= totalWeight
	sellScore /= totalWeight

	threshold := 0.6 // 可以根据需要调整这个阈值

	if buyScore > threshold {
		return &models.TradeSignal{
			StrategyName: cs.Name(),
			Action:       "BUY",
			Price:        data[len(data)-1].Close,
			Reason:       "综合策略买入信号强度高",
		}
	} else if sellScore > threshold {
		return &models.TradeSignal{
			StrategyName: cs.Name(),
			Action:       "SELL",
			Price:        data[len(data)-1].Close,
			Reason:       "综合策略卖出信号强度高",
		}
	}

	return &models.TradeSignal{
		StrategyName: cs.Name(),
		Action:       "HOLD",
		Reason:       "无明显信号",
	}
}
