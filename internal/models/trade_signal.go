package models

// TradeSignal 结构体定义了交易信号
type TradeSignal struct {
	StrategyName string
	Action       string // "BUY", "SELL", 或 "HOLD"
	Price        float64
	Reason       string
}
