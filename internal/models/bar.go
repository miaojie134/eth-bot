package models

import (
	"encoding/json"
	"time"
)

// Bar 是返回的K线数据结构
type Bar struct {
	Open       float64   `json:"open"`
	High       float64   `json:"high"`
	Low        float64   `json:"low"`
	Close      float64   `json:"close"`
	Volume     float64   `json:"volume"`
	Timestamp  time.Time `json:"timestamp"`
	TradeCount int       `json:"trade_count"`
	VWAP       float64   `json:"vwap"`
}

// AlpacaBar 是单个K线数据结构
type AlpacaBar struct {
	Open       float64 `json:"o"`
	High       float64 `json:"h"`
	Low        float64 `json:"l"`
	Close      float64 `json:"c"`
	Volume     float64 `json:"v"`
	Timestamp  string  `json:"t"`
	TradeCount int     `json:"n"`
	VWAP       float64 `json:"vw"`
}

// AlpacaBarsResponse 是 Alpaca API 的响应结构
type AlpacaBarsResponse struct {
	Bars          map[string]json.RawMessage `json:"bars"`
	NextPageToken string                     `json:"next_page_token"`
}
