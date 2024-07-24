// internal/services/analysis_service.go

package services

import (
	"time"

	"github.com/qqqq/eth-trading-system/internal/analysis"
	"github.com/qqqq/eth-trading-system/internal/analysis/indicators"
	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/storage"
)

type AnalysisService struct {
	engine   analysis.Engine
	dataRepo storage.DataRepository
}

func NewAnalysisService(dataRepo storage.DataRepository) *AnalysisService {
	analysisEngine := analysis.NewAnalysisEngine()
	analysisEngine.AddIndicator(indicators.NewSimpleMovingAverage(14))
	analysisEngine.AddIndicator(indicators.NewMACD(12, 26, 9))
	analysisEngine.AddIndicator(indicators.NewBollingerBands(20, 2))
	analysisEngine.AddIndicator(indicators.NewAverageTrueRange(14))

	return &AnalysisService{
		engine:   analysisEngine,
		dataRepo: dataRepo,
	}
}

func (s *AnalysisService) AnalyzeMarket(timeframe string, bars []models.Bar) (*analysis.AnalysisResult, error) {
	return s.engine.Analyze(bars)
}

func (s *AnalysisService) GetLatestAnalysis(timeframe string) (*analysis.AnalysisResult, error) {
	// 从数据库获取最近的一定数量的K线数据
	end := time.Now()
	start := end.Add(-24 * 30 * time.Hour) // 假设我们分析最近30天的数据
	bars, err := s.dataRepo.GetHistoricalData(timeframe, start, end)
	if err != nil {
		return nil, err
	}

	return s.AnalyzeMarket(timeframe, bars)
}
