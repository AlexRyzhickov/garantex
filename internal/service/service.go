package service

import (
	"encoding/json"
	"fmt"
	"garantex/internal/models"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Repository interface {
	Upsert(price models.Price) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) doRequest() (*models.PriceDepth, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	url := fmt.Sprintf("https://garantex.org/api/v2/depth?market=%s", "usdtrub")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	t := &models.PriceDepth{}
	err = json.Unmarshal(b, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetPrice() (models.Price, error) {
	priceDepth, err := s.doRequest()
	if err != nil {
		return models.Price{}, err
	}
	if len(priceDepth.Asks) < 1 {
		return models.Price{}, fmt.Errorf("zero len priceDepth.Asks")
	}
	if len(priceDepth.Bids) < 1 {
		return models.Price{}, fmt.Errorf("zero len priceDepth.Bids")
	}
	askPrice, err := strconv.ParseFloat(priceDepth.Asks[0].Price, 32)
	if err != nil {
		return models.Price{}, err
	}
	bidPrice, err := strconv.ParseFloat(priceDepth.Asks[0].Price, 32)
	if err != nil {
		return models.Price{}, err
	}
	price := models.Price{
		Timestamp: priceDepth.Timestamp,
		AskPrice:  askPrice,
		BidPrice:  bidPrice,
	}
	return price, nil
}

func (s *Service) SavePrice(price models.Price) error {
	return s.repo.Upsert(price)
}
