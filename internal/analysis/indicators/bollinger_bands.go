package indicators

import (
	"fmt"
	"math"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type BollingerBands struct {
	Period  int
	StdDevs float64
}

func NewBollingerBands(period int, stdDevs float64) *BollingerBands {
	return &BollingerBands{
		Period:  period,
		StdDevs: stdDevs,
	}
}

func (bb *BollingerBands) Calculate(bars []models.Bar) (interface{}, error) {
	if len(bars) < bb.Period {
		return nil, fmt.Errorf("not enough data for Bollinger Bands calculation")
	}

	sma := NewSimpleMovingAverage(bb.Period)
	middle, err := sma.Calculate(bars)
	if err != nil {
		return nil, err
	}

	middleValue := middle.(float64)

	variance := 0.0
	for i := len(bars) - bb.Period; i < len(bars); i++ {
		variance += math.Pow(bars[i].Close-middleValue, 2)
	}
	stdDev := math.Sqrt(variance / float64(bb.Period))

	upper := middleValue + bb.StdDevs*stdDev
	lower := middleValue - bb.StdDevs*stdDev

	return map[string]float64{
		"Upper":  upper,
		"Middle": middleValue,
		"Lower":  lower,
	}, nil
}

func (bb *BollingerBands) Name() string {
	return fmt.Sprintf("BB(%d,%.1f)", bb.Period, bb.StdDevs)
}
