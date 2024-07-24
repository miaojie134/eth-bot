package api

import (
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/services"
)

type Handler struct {
	alpacaService         *services.AlpacaService
	dataCollectionService *services.DataCollectionService
}

func NewHandler(alpacaService *services.AlpacaService, dataCollectionService *services.DataCollectionService) *Handler {
	return &Handler{
		alpacaService:         alpacaService,
		dataCollectionService: dataCollectionService,
	}
}

func (h *Handler) GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	bar, err := h.dataCollectionService.GetLatestPrice()
	if err != nil {
		http.Error(w, "获取最新价格失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"open":       bar.Open,
		"high":       bar.High,
		"low":        bar.Low,
		"close":      bar.Close,
		"volume":     bar.Volume,
		"timestamp":  bar.Timestamp,
		"tradeCount": bar.TradeCount,
		"vwap":       bar.VWAP,
	})
}

func (h *Handler) GetHistoricalData(w http.ResponseWriter, r *http.Request) {
	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		http.Error(w, "必须提供timeframe参数", http.StatusBadRequest)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "开始时间无效", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "结束时间无效", http.StatusBadRequest)
		return
	}

	bars, err := h.dataCollectionService.GetHistoricalData(timeframe, start, end)
	if err != nil {
		http.Error(w, "获取历史数据失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bars)
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	timeframe := "1Day"
	latestBar, err := h.dataCollectionService.GetLatestPrice()
	if err != nil {
		http.Error(w, "获取最新数据失败", http.StatusInternalServerError)
		return
	}

	end := time.Now()
	start := end.Add(-24 * time.Hour)
	historicalData, err := h.dataCollectionService.GetHistoricalData(timeframe, start, end)
	if err != nil {
		http.Error(w, "获取历史数据失败", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsRunning             bool
		LatestPrice           float64
		HistoricalData        []models.Bar
		DataCollectionRunning bool
	}{
		IsRunning:             true,
		LatestPrice:           latestBar.Close,
		HistoricalData:        historicalData,
		DataCollectionRunning: true,
	}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "加载模板失败", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "渲染模板失败", http.StatusInternalServerError)
		return
	}
}
