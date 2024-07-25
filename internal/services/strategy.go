package services

import (
	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/strategy"
)

type StrategyService struct {
	strategies []strategy.Strategy
}

func NewStrategyService(strategies ...strategy.Strategy) *StrategyService {
	return &StrategyService{
		strategies: strategies,
	}
}

func (s *StrategyService) AddStrategy(strategy strategy.Strategy) {
	s.strategies = append(s.strategies, strategy)
}

func (s *StrategyService) EvaluateStrategies(data []models.Bar, analysisResult *models.AnalysisResult) []*models.TradeSignal {
	var signals []*models.TradeSignal
	for _, strategy := range s.strategies {
		signal := strategy.Evaluate(data, analysisResult)
		signals = append(signals, signal)
	}
	return signals
}
