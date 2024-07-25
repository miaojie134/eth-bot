// internal/services/analysis_service.go
package services

import (
	"time"

	"github.com/qqqq/eth-trading-system/internal/analysis"
	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/storage"
)

type AnalysisService struct {
	engine          analysis.Engine
	dataRepo        storage.DataRepository
	strategyService *StrategyService
}

func NewAnalysisService(engine analysis.Engine, dataRepo storage.DataRepository, strategyService *StrategyService) *AnalysisService {
	return &AnalysisService{
		engine:          engine,
		dataRepo:        dataRepo,
		strategyService: strategyService,
	}
}

func (s *AnalysisService) AnalyzeMarket(timeframe string, bars []models.Bar) (*models.AnalysisResult, error) {
	analysisResult, err := s.engine.Analyze(bars)
	if err != nil {
		return nil, err
	}

	signals := s.strategyService.EvaluateStrategies(bars, analysisResult)
	analysisResult.StrategySignals = signals

	return analysisResult, nil
}

func (s *AnalysisService) GetLatestAnalysis(timeframe string) (*models.AnalysisResult, error) {
	// 从数据库获取最近的一定数量的K线数据
	end := time.Now()
	start := end.Add(-24 * 30 * time.Hour) // 假设我们分析最近30天的数据
	bars, err := s.dataRepo.GetHistoricalData(timeframe, start, end)
	if err != nil {
		return nil, err
	}

	return s.AnalyzeMarket(timeframe, bars)
}

// 添加新方法以获取特定策略的分析结果
func (s *AnalysisService) GetStrategyAnalysis(timeframe string, strategyName string) (*models.TradeSignal, error) {
	result, err := s.GetLatestAnalysis(timeframe)
	if err != nil {
		return nil, err
	}

	for _, signal := range result.StrategySignals {
		if signal.StrategyName == strategyName {
			return signal, nil
		}
	}

	return nil, nil // 如果没有找到指定策略的信号，返回 nil
}
