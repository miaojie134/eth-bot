package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/qqqq/eth-trading-system/internal/datamanager"
	"github.com/qqqq/eth-trading-system/internal/models"
	"github.com/qqqq/eth-trading-system/internal/utils"
)

type AlpacaClient interface {
	GetHistoricalBars(symbol, timeframe, start, end string, limit int, pageToken string) ([]models.AlpacaBar, string, error)
	GetLatestBar(symbol string) (*models.AlpacaBar, error)
}

type HTTPAlpacaClient struct {
	apiKey    string
	apiSecret string
	client    *http.Client
	baseURL   string
}

func NewAlpacaClient(apiKey, apiSecret string) AlpacaClient {
	return &HTTPAlpacaClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   "https://data.alpaca.markets/v1beta3/crypto",
	}
}

func (c *HTTPAlpacaClient) GetHistoricalBars(symbol, timeframe, start, end string, limit int, pageToken string) ([]models.AlpacaBar, string, error) {
	endpoint := fmt.Sprintf("%s/us/bars", c.baseURL)

	params := url.Values{}
	params.Add("symbols", symbol)
	params.Add("timeframe", timeframe)
	if start != "" {
		params.Add("start", start)
	}
	if end != "" {
		params.Add("end", end)
	}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	if pageToken != "" {
		params.Add("page_token", pageToken)
	}

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, "", err
	}

	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var response models.AlpacaBarsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, "", err
	}

	rawBars, ok := response.Bars[symbol]
	if !ok {
		return nil, "", fmt.Errorf("没有数据可用: %s", symbol)
	}

	var bars []models.AlpacaBar
	err = json.Unmarshal(rawBars, &bars)
	if err != nil {
		return nil, "", err
	}

	return bars, response.NextPageToken, nil
}

func (c *HTTPAlpacaClient) GetLatestBar(symbol string) (*models.AlpacaBar, error) {
	endpoint := fmt.Sprintf("%s/us/latest/bars", c.baseURL)

	params := url.Values{}
	params.Add("symbols", symbol)

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	utils.Log.Infof("请求最新数据: %s?%s", endpoint, params.Encode())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response models.AlpacaBarsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	rawBar, ok := response.Bars[symbol]
	if !ok || len(rawBar) == 0 {
		return nil, fmt.Errorf("没有数据可用: %s", symbol)
	}

	var bar models.AlpacaBar
	err = json.Unmarshal(rawBar, &bar)
	if err != nil {
		return nil, err
	}

	return &bar, nil
}

func (c *HTTPAlpacaClient) setHeaders(req *http.Request) {
	req.Header.Set("APCA-API-KEY-ID", c.apiKey)
	req.Header.Set("APCA-API-SECRET-KEY", c.apiSecret)
	req.Header.Add("accept", "application/json")
}

type AlpacaService struct {
	client AlpacaClient
}

// 确保 AlpacaService 实现了 MarketDataProvider 接口
var _ datamanager.MarketDataProvider = (*AlpacaService)(nil)

func NewAlpacaService(client AlpacaClient) *AlpacaService {
	return &AlpacaService{client: client}
}

func (s *AlpacaService) GetHistoricalBars(symbol, timeframe, start, end string, limit int, pageToken string) ([]models.AlpacaBar, string, error) {
	return s.client.GetHistoricalBars(symbol, timeframe, start, end, limit, pageToken)
}

func (s *AlpacaService) GetLatestBar(symbol string) (*models.AlpacaBar, error) {
	return s.client.GetLatestBar(symbol)
}
