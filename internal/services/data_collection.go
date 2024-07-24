package services

import (
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/storage"
	"github.com/qqqq/eth-trading-system/internal/utils"
)

const (
	ethSymbol = "ETH/USD"
	apiLimit  = 10000
)

type DataCollectionService struct {
	alpacaService *AlpacaService
	dataRepo      storage.DataRepository
}

func NewDataCollectionService(alpacaService *AlpacaService, dataRepo storage.DataRepository) *DataCollectionService {
	return &DataCollectionService{
		alpacaService: alpacaService,
		dataRepo:      dataRepo,
	}
}

func (s *DataCollectionService) Start() {
	s.initializeData()

	go s.startTicker("5Min", 5*time.Minute)
	go s.startTicker("15Min", 15*time.Minute)
	go s.startTicker("1Hour", 1*time.Hour)
	go s.startTicker("4Hour", 4*time.Hour)
	go s.startTicker("1Day", 24*time.Hour)
	utils.Log.Info("数据收集服务已启动")
}

func (s *DataCollectionService) initializeData() {
	s.initializeHistoricalData()
	s.collectAndStoreLatestPrice()
}

func (s *DataCollectionService) initializeHistoricalData() {
	timeframes := []string{"5Min", "15Min", "1Hour", "4Hour", "1Day"}
	for _, timeframe := range timeframes {
		utils.Log.Infof("初始化历史数据，时间框架: %s", timeframe)
		s.collectAndStoreHistoricalData(timeframe)
	}
}

func (s *DataCollectionService) startTicker(timeframe string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		utils.Log.Infof("收集并存储数据，时间框架: %s", timeframe)
		s.collectAndStoreHistoricalData(timeframe)
	}
}

func (s *DataCollectionService) collectAndStoreLatestPrice() {
	bar, err := s.alpacaService.GetLatestBar(ethSymbol)
	if err != nil {
		utils.Log.WithError(err).Error("获取最新数据失败")
		return
	}

	latestBar := convertAlpacaBarToBar(bar)
	err = s.dataRepo.StoreLatestPrice(latestBar)
	if err != nil {
		utils.Log.WithError(err).Error("存储最新数据失败")
		return
	}

	utils.Log.WithFields(map[string]interface{}{
		"symbol": ethSymbol,
		"price":  latestBar.Close,
		"time":   latestBar.Timestamp,
	}).Info("最新价格已收集并存储")
}

func (s *DataCollectionService) collectAndStoreHistoricalData(timeframe string) {
	end := time.Now().UTC()
	start := end.Add(-calculateDuration(timeframe))
	utils.Log.Infof("收集历史数据，时间框架: %s, 开始: %s, 结束: %s", timeframe, start, end)

	pageToken := ""
	pageCount := 0

	for {
		pageCount++
		bars, newPageToken, err := s.alpacaService.GetHistoricalBars(ethSymbol, timeframe, start.Format(time.RFC3339), end.Format(time.RFC3339), apiLimit, pageToken)
		if err != nil {
			utils.Log.WithError(err).Error("获取历史数据失败")
			return
		}

		convertedBars := convertAlpacaBarsToBar(bars)
		err = s.dataRepo.StoreHistoricalData(convertedBars, timeframe)
		if err != nil {
			utils.Log.WithError(err).Error("存储历史数据失败")
			return
		}

		if newPageToken == "" {
			utils.Log.Infof("历史数据收集完成，时间框架: %s, 页数: %d", timeframe, pageCount)
			break
		}

		pageToken = newPageToken
		utils.Log.Infof("历史数据收集中，时间框架: %s, 页数: %d", timeframe, pageCount)
	}
}

func (s *DataCollectionService) GetLatestPrice() (*models.Bar, error) {
	return s.dataRepo.GetLatestPrice()
}

func (s *DataCollectionService) GetHistoricalData(timeframe string, start, end time.Time) ([]models.Bar, error) {
	return s.dataRepo.GetHistoricalData(timeframe, start, end)
}

func calculateDuration(timeframe string) time.Duration {
	switch timeframe {
	case "5Min":
		return 2 * 7 * 24 * time.Hour // 2 weeks
	case "15Min":
		return 30 * 24 * time.Hour // 1 month
	case "1Hour":
		return 3 * 30 * 24 * time.Hour // 3 months
	case "4Hour":
		return 12 * 30 * 24 * time.Hour // 1 year
	case "1Day":
		return 2 * 365 * 24 * time.Hour // 2 years
	default:
		return 0
	}
}

func convertAlpacaBarToBar(alpacaBar *models.AlpacaBar) *models.Bar {
	timestamp, _ := time.Parse(time.RFC3339, alpacaBar.Timestamp)
	return &models.Bar{
		Open:       alpacaBar.Open,
		High:       alpacaBar.High,
		Low:        alpacaBar.Low,
		Close:      alpacaBar.Close,
		Volume:     alpacaBar.Volume,
		Timestamp:  timestamp,
		TradeCount: alpacaBar.TradeCount,
		VWAP:       alpacaBar.VWAP,
	}
}

func convertAlpacaBarsToBar(alpacaBars []models.AlpacaBar) []models.Bar {
	var bars []models.Bar
	for _, alpacaBar := range alpacaBars {
		bars = append(bars, *convertAlpacaBarToBar(&alpacaBar))
	}
	return bars
}
