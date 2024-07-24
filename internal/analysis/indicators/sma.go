package indicators

import (
	"fmt"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type SimpleMovingAverage struct {
	Period int
}

func NewSimpleMovingAverage(period int) *SimpleMovingAverage {
	return &SimpleMovingAverage{Period: period}
}

func (sma *SimpleMovingAverage) Calculate(bars []models.Bar) (interface{}, error) {
	if len(bars) < sma.Period {
		return nil, fmt.Errorf("not enough data for SMA calculation")
	}
	sum := 0.0
	for i := len(bars) - sma.Period; i < len(bars); i++ {
		sum += bars[i].Close
	}
	return sum / float64(sma.Period), nil
}

func (sma *SimpleMovingAverage) Name() string {
	return fmt.Sprintf("SMA%d", sma.Period)
}
