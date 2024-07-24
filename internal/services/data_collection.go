package services

import (
	"database/sql"
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/utils"
	"github.com/sirupsen/logrus"
)

const apiLimit = 10000 // API的限制条数

type DataCollectionService struct {
	alpacaService *AlpacaService
	db            *sql.DB
}

func NewDataCollectionService(alpacaService *AlpacaService, db *sql.DB) *DataCollectionService {
	return &DataCollectionService{
		alpacaService: alpacaService,
		db:            db,
	}
}

func (s *DataCollectionService) Start() {
	// 初始化时先获取历史数据和最新价格
	s.initializeData()

	// 设置不同K线周期的更新频率
	go s.startTicker("5Min", 5*time.Minute)
	go s.startTicker("15Min", 15*time.Minute)
	go s.startTicker("1Hour", 1*time.Hour)
	go s.startTicker("4Hour", 4*time.Hour)
	go s.startTicker("1Day", 24*time.Hour)
	utils.Log.Info("Data collection service started")
}

func (s *DataCollectionService) initializeData() {
	s.initializeHistoricalData()
	s.collectAndStoreLatestPrice()
}

func (s *DataCollectionService) initializeHistoricalData() {
	timeframes := []string{"5Min", "15Min", "1Hour", "4Hour", "1Day"}
	for _, timeframe := range timeframes {
		utils.Log.Infof("初始化历史数据，周期: %s", timeframe)
		s.collectAndStoreHistoricalData(timeframe)
	}
}

func (s *DataCollectionService) startTicker(timeframe string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		utils.Log.Infof("收集和存储数据，周期: %s", timeframe)
		s.collectAndStoreHistoricalData(timeframe)
	}
}

func (s *DataCollectionService) collectAndStoreLatestPrice() {
	bar, err := s.alpacaService.GetLatestBar()
	if err != nil {
		utils.Log.WithError(err).Error("Failed to get latest data")
		return
	}

	timestamp, err := time.Parse(time.RFC3339, bar.Timestamp)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to parse timestamp")
		return
	}

	_, err = s.db.Exec(`
		INSERT INTO latest_price (symbol, open, high, low, close, volume, timestamp, trade_count, vwap)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(symbol) DO UPDATE SET
		open=excluded.open, high=excluded.high, low=excluded.low, close=excluded.close, volume=excluded.volume, timestamp=excluded.timestamp, trade_count=excluded.trade_count, vwap=excluded.vwap
	`, "ETH/USD", bar.Open, bar.High, bar.Low, bar.Close, bar.Volume, timestamp, bar.TradeCount, bar.VWAP)

	if err != nil {
		utils.Log.WithError(err).Error("Failed to store latest data")
		return
	}
	utils.Log.WithFields(logrus.Fields{
		"symbol": "ETH/USD",
		"price":  bar.Close,
		"time":   bar.Timestamp,
	}).Info("Latest price collected and stored")

}

func (s *DataCollectionService) collectAndStoreHistoricalData(timeframe string) {
	end := time.Now().UTC()
	start := end.Add(-calculateDuration(timeframe))
	utils.Log.Infof("收集历史数据，周期: %s, 开始时间: %s, 结束时间: %s", timeframe, start, end)
	pageToken := ""
	pageCount := 0
	for {
		pageCount++
		bars, newPageToken, err := s.alpacaService.GetHistoricalBars("ETH/USD", timeframe, start.Format(time.RFC3339), end.Format(time.RFC3339), apiLimit, pageToken)
		if err != nil {
			utils.Log.WithError(err).Error("Failed to get historical data")
			return
		}

		for _, bar := range bars {
			timestamp, err := time.Parse(time.RFC3339, bar.Timestamp)
			if err != nil {
				utils.Log.WithError(err).Error("Failed to parse timestamp")
				continue
			}

			_, err = s.db.Exec(`
				INSERT INTO market_data (symbol, timeframe, open, high, low, close, volume, timestamp, trade_count, vwap)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(symbol, timeframe, timestamp) DO UPDATE SET
				open=excluded.open, high=excluded.high, low=excluded.low, close=excluded.close, volume=excluded.volume, trade_count=excluded.trade_count, vwap=excluded.vwap
			`, "ETH/USD", timeframe, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume, timestamp, bar.TradeCount, bar.VWAP)

			if err != nil {
				utils.Log.WithError(err).Error("Failed to store historical data")
				return
			}
		}

		if newPageToken == "" {
			utils.Log.Infof("历史数据收集完成，周期: %s，页数: %d", timeframe, pageCount)
			break
		}

		pageToken = newPageToken
		utils.Log.Infof("历史数据收集完成，周期: %s，页数: %d", timeframe, pageCount)
	}
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

func (s *DataCollectionService) GetLatestStoredData(timeframe string) (*models.Bar, error) {
	var bar models.Bar
	err := s.db.QueryRow(`
		SELECT open, high, low, close, volume, timestamp, trade_count, vwap
		FROM market_data
		WHERE symbol = 'ETH/USD' AND timeframe = ?
		ORDER BY timestamp DESC
		LIMIT 1
	`, timeframe).Scan(&bar.Open, &bar.High, &bar.Low, &bar.Close, &bar.Volume, &bar.Timestamp, &bar.TradeCount, &bar.VWAP)

	if err != nil {
		return nil, err
	}

	return &bar, nil
}

func (s *DataCollectionService) GetLatestPrice() (*models.Bar, error) {
	var bar models.Bar
	err := s.db.QueryRow(`
		SELECT open, high, low, close, volume, timestamp, trade_count, vwap
		FROM latest_price
		WHERE symbol = 'ETH/USD'
	`).Scan(&bar.Open, &bar.High, &bar.Low, &bar.Close, &bar.Volume, &bar.Timestamp, &bar.TradeCount, &bar.VWAP)

	if err != nil {
		return nil, err
	}

	return &bar, nil
}

func (s *DataCollectionService) GetHistoricalData(timeframe string, start, end time.Time) ([]models.Bar, error) {
	rows, err := s.db.Query(`
		SELECT open, high, low, close, volume, timestamp, trade_count, vwap
		FROM market_data
		WHERE symbol = 'ETH/USD' AND timeframe = ? AND timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`, timeframe, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bars []models.Bar
	for rows.Next() {
		var bar models.Bar
		err := rows.Scan(&bar.Open, &bar.High, &bar.Low, &bar.Close, &bar.Volume, &bar.Timestamp, &bar.TradeCount, &bar.VWAP)
		if err != nil {
			return nil, err
		}
		bars = append(bars, bar)
	}

	return bars, nil
}
