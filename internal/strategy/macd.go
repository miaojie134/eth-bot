package strategy

import (
	"github.com/qqqq/eth-trading-system/internal/analysis"
	"github.com/qqqq/eth-trading-system/internal/models"
)

// MACDStrategy 实现了MACD策略
type MACDStrategy struct {
	BaseStrategy
}

func NewMACDStrategy() *MACDStrategy {
	return &MACDStrategy{
		BaseStrategy: BaseStrategy{name: "MACD"},
	}
}

func (s *MACDStrategy) Evaluate(data []models.Bar, analysisResult *models.AnalysisResult) *models.TradeSignal {
	macdLine, ok := analysisResult.Indicators[analysis.IndicatorMACD].([]float64)
	if !ok {
		return &models.TradeSignal{StrategyName: s.Name(), Action: "HOLD", Reason: "MACD指标不可用"}
	}

	signalLine, ok := analysisResult.Indicators[analysis.IndicatorMACDSignal].([]float64)
	if !ok {
		return &models.TradeSignal{StrategyName: s.Name(), Action: "HOLD", Reason: "MACD信号线不可用"}
	}

	if len(macdLine) < 2 || len(signalLine) < 2 {
		return &models.TradeSignal{StrategyName: s.Name(), Action: "HOLD", Reason: "不足够的数据"}
	}

	if macdLine[len(macdLine)-1] > signalLine[len(signalLine)-1] &&
		macdLine[len(macdLine)-2] <= signalLine[len(signalLine)-2] {
		return &models.TradeSignal{
			StrategyName: s.Name(),
			Action:       "BUY",
			Price:        data[len(data)-1].Close,
			Reason:       "MACD线上穿信号线",
		}
	} else if macdLine[len(macdLine)-1] < signalLine[len(signalLine)-1] &&
		macdLine[len(macdLine)-2] >= signalLine[len(signalLine)-2] {
		return &models.TradeSignal{
			StrategyName: s.Name(),
			Action:       "SELL",
			Price:        data[len(data)-1].Close,
			Reason:       "MACD线下穿信号线",
		}
	}

	return &models.TradeSignal{StrategyName: s.Name(), Action: "HOLD", Reason: "无显著变化"}
}
