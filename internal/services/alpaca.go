package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/qqqq/eth-trading-system/internal/models"
)

const (
	baseURL   = "https://data.alpaca.markets/v1beta3/crypto"
	ethSymbol = "ETH/USD"
)

type AlpacaService struct {
	apiKey    string
	apiSecret string
	client    *http.Client
}

func NewAlpacaService(apiKey, apiSecret string) *AlpacaService {
	return &AlpacaService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *AlpacaService) GetHistoricalBars(symbol, timeframe, start, end string, limit int, pageToken string) ([]models.AlpacaBar, string, error) {
	endpoint := fmt.Sprintf("%s/us/bars", baseURL)

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

	s.setHeaders(req)

	resp, err := s.client.Do(req)
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
		return nil, "", fmt.Errorf("没有可用数据: %s", symbol)
	}

	var bars []models.AlpacaBar
	err = json.Unmarshal(rawBars, &bars)
	if err != nil {
		return nil, "", err
	}

	return bars, response.NextPageToken, nil
}

func (s *AlpacaService) GetLatestBar() (*models.AlpacaBar, error) {
	endpoint := fmt.Sprintf("%s/us/latest/bars", baseURL)

	params := url.Values{}
	params.Add("symbols", ethSymbol)

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	s.setHeaders(req)

	log.Printf("请求最新数据: %s?%s", endpoint, params.Encode())

	resp, err := s.client.Do(req)
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

	rawBar, ok := response.Bars[ethSymbol]
	if !ok || len(rawBar) == 0 {
		return nil, fmt.Errorf("没有可用数据: ETH/USD")
	}

	var bar models.AlpacaBar
	err = json.Unmarshal(rawBar, &bar)
	if err != nil {
		return nil, err
	}

	return &bar, nil
}

func (s *AlpacaService) setHeaders(req *http.Request) {
	req.Header.Set("APCA-API-KEY-ID", s.apiKey)
	req.Header.Set("APCA-API-SECRET-KEY", s.apiSecret)
	req.Header.Add("accept", "application/json")
}
