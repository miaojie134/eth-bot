package api

import (
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"github.com/qqqq/eth-trading-system/internal/analysis"
	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/services"
)

type Handler struct {
	alpacaService         *services.AlpacaService
	dataCollectionService *services.DataCollectionService
	analysisService       *services.AnalysisService
}

func NewHandler(alpacaService *services.AlpacaService, dataCollectionService *services.DataCollectionService, analysisService *services.AnalysisService) *Handler {
	return &Handler{
		alpacaService:         alpacaService,
		dataCollectionService: dataCollectionService,
		analysisService:       analysisService,
	}
}

func (h *Handler) GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	bar, err := h.dataCollectionService.GetLatestPrice()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "获取最新价格失败")
		return
	}
	respondWithJSON(w, http.StatusOK, bar)
}

func (h *Handler) GetHistoricalData(w http.ResponseWriter, r *http.Request) {
	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		respondWithError(w, http.StatusBadRequest, "必须提供timeframe参数")
		return
	}

	start, end, err := parseTimeRange(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bars, err := h.dataCollectionService.GetHistoricalData(timeframe, start, end)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "获取历史数据失败")
		return
	}
	respondWithJSON(w, http.StatusOK, bars)
}

func (h *Handler) GetMarketAnalysis(w http.ResponseWriter, r *http.Request) {
	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		respondWithError(w, http.StatusBadRequest, "必须提供timeframe参数")
		return
	}

	result, err := h.analysisService.GetLatestAnalysis(timeframe)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "获取市场分析失败")
		return
	}
	respondWithJSON(w, http.StatusOK, result)
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := h.getIndexData()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "获取数据失败")
		return
	}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "加载模板失败")
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "渲染模板失败")
	}
}

// Helper functions

func (h *Handler) getIndexData() (interface{}, error) {
	latestBar, err := h.dataCollectionService.GetLatestPrice()
	if err != nil {
		return nil, err
	}

	analysisResult, err := h.analysisService.GetLatestAnalysis("1Day")
	if err != nil {
		return nil, err
	}

	end := time.Now()
	start := end.Add(-24 * time.Hour)
	historicalData, err := h.dataCollectionService.GetHistoricalData("1Day", start, end)
	if err != nil {
		return nil, err
	}

	return struct {
		IsRunning             bool
		LatestPrice           float64
		HistoricalData        []models.Bar
		DataCollectionRunning bool
		MarketAnalysis        *analysis.AnalysisResult
	}{
		IsRunning:             true,
		LatestPrice:           latestBar.Close,
		HistoricalData:        historicalData,
		DataCollectionRunning: true,
		MarketAnalysis:        analysisResult,
	}, nil
}

func parseTimeRange(r *http.Request) (time.Time, time.Time, error) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
