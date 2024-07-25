package models

import "time"

type AnalysisResult struct {
	Timestamp       time.Time
	MarketState     MarketState
	Indicators      map[string]interface{}
	StrategySignals []*TradeSignal
}

type MarketState int

const (
	Bullish MarketState = iota
	Bearish
	Neutral
)
