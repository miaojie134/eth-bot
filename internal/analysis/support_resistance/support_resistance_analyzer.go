package support_resistance

import (
	"sort"

	"github.com/qqqq/eth-trading-system/internal/models"
)

type SupportResistanceAnalyzer struct {
	WindowSize int
}

func NewSupportResistanceAnalyzer(windowSize int) *SupportResistanceAnalyzer {
	return &SupportResistanceAnalyzer{WindowSize: windowSize}
}

func (sra *SupportResistanceAnalyzer) FindLevels(bars []models.Bar) ([]float64, []float64) {
	if len(bars) < sra.WindowSize*2 {
		return nil, nil
	}

	var supports, resistances []float64

	for i := sra.WindowSize; i < len(bars)-sra.WindowSize; i++ {
		if isSupport(bars, i, sra.WindowSize) {
			supports = append(supports, bars[i].Low)
		}
		if isResistance(bars, i, sra.WindowSize) {
			resistances = append(resistances, bars[i].High)
		}
	}

	supports = removeDuplicates(supports)
	resistances = removeDuplicates(resistances)

	sort.Float64s(supports)
	sort.Float64s(resistances)

	return supports, resistances
}

func isSupport(bars []models.Bar, index, windowSize int) bool {
	for i := index - windowSize; i < index; i++ {
		if bars[i].Low <= bars[index].Low {
			return false
		}
	}
	for i := index + 1; i <= index+windowSize; i++ {
		if bars[i].Low <= bars[index].Low {
			return false
		}
	}
	return true
}

func isResistance(bars []models.Bar, index, windowSize int) bool {
	for i := index - windowSize; i < index; i++ {
		if bars[i].High >= bars[index].High {
			return false
		}
	}
	for i := index + 1; i <= index+windowSize; i++ {
		if bars[i].High >= bars[index].High {
			return false
		}
	}
	return true
}

func removeDuplicates(slice []float64) []float64 {
	keys := make(map[float64]bool)
	var list []float64
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
