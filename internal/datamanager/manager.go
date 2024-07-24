// internal/datamanager/manager.go
package datamanager

import (
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/storage"
	"github.com/qqqq/eth-trading-system/internal/utils"
)

const (
	EthSymbol = "ETH/USD"
	ApiLimit  = 10000
)

// MarketDataProvider 定义了获取市场数据的接口
type MarketDataProvider interface {
	GetLatestBar(symbol string) (*models.AlpacaBar, error)
	GetHistoricalBars(symbol, timeframe, start, end string, limit int, pageToken string) ([]models.AlpacaBar, string, error)
}

type DataManager struct {
	marketDataProvider MarketDataProvider
	dataRepo           storage.DataRepository
}

func NewDataManager(marketDataProvider MarketDataProvider, dataRepo storage.DataRepository) *DataManager {
	return &DataManager{
		marketDataProvider: marketDataProvider,
		dataRepo:           dataRepo,
	}
}

func (dm *DataManager) CollectAndStoreLatestPrice() error {
	bar, err := dm.marketDataProvider.GetLatestBar(EthSymbol)
	if err != nil {
		return err
	}

	latestBar := convertAlpacaBarToBar(bar)
	err = dm.dataRepo.StoreLatestPrice(latestBar)
	if err != nil {
		return err
	}

	utils.Log.WithFields(map[string]interface{}{
		"symbol": EthSymbol,
		"price":  latestBar.Close,
		"time":   latestBar.Timestamp,
	}).Info("最新价格已收集并存储")

	return nil
}

func (dm *DataManager) CollectAndStoreHistoricalData(timeframe string, start, end time.Time) error {
	pageToken := ""
	pageCount := 0

	for {
		pageCount++
		bars, newPageToken, err := dm.marketDataProvider.GetHistoricalBars(EthSymbol, timeframe, start.Format(time.RFC3339), end.Format(time.RFC3339), ApiLimit, pageToken)
		if err != nil {
			return err
		}

		convertedBars := convertAlpacaBarsToBar(bars)
		err = dm.dataRepo.StoreHistoricalData(convertedBars, timeframe)
		if err != nil {
			return err
		}

		if newPageToken == "" {
			utils.Log.Infof("历史数据收集完成，时间框架: %s, 页数: %d", timeframe, pageCount)
			break
		}

		pageToken = newPageToken
		utils.Log.Infof("历史数据收集中，时间框架: %s, 页数: %d", timeframe, pageCount)
	}

	return nil
}

func (dm *DataManager) GetLatestPrice() (*models.Bar, error) {
	return dm.dataRepo.GetLatestPrice()
}

func (dm *DataManager) GetHistoricalData(timeframe string, start, end time.Time) ([]models.Bar, error) {
	return dm.dataRepo.GetHistoricalData(timeframe, start, end)
}

// 辅助函数
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
