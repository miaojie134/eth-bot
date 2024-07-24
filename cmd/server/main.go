package main

import (
	"log"
	"net/http"

	"github.com/qqqq/eth-trading-system/internal/api"
	"github.com/qqqq/eth-trading-system/internal/config"
	"github.com/qqqq/eth-trading-system/internal/services"
	"github.com/qqqq/eth-trading-system/internal/storage"
	"github.com/qqqq/eth-trading-system/internal/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	if err := utils.InitLogger(cfg.LogDir, cfg.LogLevel); err != nil {
		log.Fatalf("初始化日志记录器失败: %v", err)
	}

	utils.Log.Info("启动ETH交易系统")

	db, err := storage.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		utils.Log.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close()

	dataRepo := storage.NewDataRepository(db.DB)
	analysisService := services.NewAnalysisService(dataRepo)

	alpacaClient := services.NewAlpacaClient(cfg.AlpacaAPIKey, cfg.AlpacaAPISecret)
	alpacaService := services.NewAlpacaService(alpacaClient)
	dataCollectionService := services.NewDataCollectionService(alpacaService, dataRepo)
	dataCollectionService.Start()

	handler := api.NewHandler(alpacaService, dataCollectionService, analysisService)

	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/api/price", handler.GetLatestPrice)
	http.HandleFunc("/api/historical", handler.GetHistoricalData)
	http.HandleFunc("/api/analysis", handler.GetMarketAnalysis)

	utils.Log.Infof("服务器启动在端口 %s", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, nil); err != nil {
		utils.Log.Fatalf("启动服务器失败: %v", err)
	}
}
