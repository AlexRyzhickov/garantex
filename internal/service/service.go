package service

import (
	"fmt"
	"garantex/internal/client"
	"garantex/internal/models"
	"strconv"
)

type Repository interface {
	Upsert(price models.Price) error
}

type Client interface {
	DoRequest() (*models.PriceDepth, error)
}

type Service struct {
	repo   Repository
	client Client
}

func New(repo Repository) *Service {
	return &Service{
		repo:   repo,
		client: &client.ApiClient{},
	}
}

func (s *Service) GetPrice() (models.Price, error) {
	priceDepth, err := s.client.DoRequest()
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
	bidPrice, err := strconv.ParseFloat(priceDepth.Bids[0].Price, 32)
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
