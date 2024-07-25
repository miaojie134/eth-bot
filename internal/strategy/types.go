package strategy

import (
	"github.com/qqqq/eth-trading-system/internal/models"
)

// Strategy 接口定义了策略的基本行为
type Strategy interface {
	Evaluate(data []models.Bar, analysisResult *models.AnalysisResult) *models.TradeSignal
	Name() string
}

// BaseStrategy 提供了基础的策略功能
type BaseStrategy struct {
	name string
}

func (bs *BaseStrategy) Name() string {
	return bs.name
}
