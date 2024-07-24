package indicators

import (
	"fmt"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type RelativeStrengthIndex struct {
	Period int
}

func NewRelativeStrengthIndex(period int) *RelativeStrengthIndex {
	return &RelativeStrengthIndex{
		Period: period,
	}
}

func (rsi *RelativeStrengthIndex) Calculate(bars []models.Bar) (interface{}, error) {
	if len(bars) < rsi.Period+1 {
		return nil, fmt.Errorf("not enough data for RSI calculation, need at least %d bars", rsi.Period+1)
	}

	var gains, losses float64
	for i := 1; i <= rsi.Period; i++ {
		change := bars[i].Close - bars[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(rsi.Period)
	avgLoss := losses / float64(rsi.Period)

	for i := rsi.Period + 1; i < len(bars); i++ {
		change := bars[i].Close - bars[i-1].Close
		if change > 0 {
			avgGain = (avgGain*float64(rsi.Period-1) + change) / float64(rsi.Period)
			avgLoss = (avgLoss*float64(rsi.Period-1) + 0) / float64(rsi.Period)
		} else {
			avgGain = (avgGain*float64(rsi.Period-1) + 0) / float64(rsi.Period)
			avgLoss = (avgLoss*float64(rsi.Period-1) - change) / float64(rsi.Period)
		}
	}

	if avgLoss == 0 {
		return 100.0, nil
	}

	rs := avgGain / avgLoss
	rsiValue := 100 - (100 / (1 + rs))

	return rsiValue, nil
}

func (rsi *RelativeStrengthIndex) Name() string {
	return fmt.Sprintf("RSI(%d)", rsi.Period)
}
