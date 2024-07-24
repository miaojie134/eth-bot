package indicators

import (
	"fmt"
	"math"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type AverageTrueRange struct {
	Period int
}

func NewAverageTrueRange(period int) *AverageTrueRange {
	return &AverageTrueRange{Period: period}
}

func (atr *AverageTrueRange) Calculate(bars []models.Bar) (interface{}, error) {
	if len(bars) < atr.Period+1 {
		return nil, fmt.Errorf("not enough data for ATR calculation")
	}

	trueRanges := make([]float64, len(bars)-1)
	for i := 1; i < len(bars); i++ {
		high := bars[i].High
		low := bars[i].Low
		prevClose := bars[i-1].Close

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)

		trueRanges[i-1] = math.Max(tr1, math.Max(tr2, tr3))
	}

	atrValue := 0.0
	for i := 0; i < atr.Period; i++ {
		atrValue += trueRanges[i]
	}
	atrValue /= float64(atr.Period)

	for i := atr.Period; i < len(trueRanges); i++ {
		atrValue = (atrValue*(float64(atr.Period)-1) + trueRanges[i]) / float64(atr.Period)
	}

	return atrValue, nil
}

func (atr *AverageTrueRange) Name() string {
	return fmt.Sprintf("ATR(%d)", atr.Period)
}
