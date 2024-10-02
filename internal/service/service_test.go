package service

import (
	"errors"
	"garantex/internal/mock"
	"garantex/internal/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestService(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	client  mock.Client
	service Service
}

func (s *Suite) SetupTest() {
	s.client = mock.Client{}
	s.service = Service{
		client: &s.client,
	}
}

func (s *Suite) TestMajor() {
	s.client.On("DoRequest").Return(&models.PriceDepth{
		Asks: []models.Ask{{Price: "67.15"}},
		Bids: []models.Bid{{Price: "67.15"}},
	}, nil)
	_, err := s.service.GetPrice()
	s.NoError(err)
}

func (s *Suite) TestMinorClientReturnErr() {
	s.client.On("DoRequest").Return(&models.PriceDepth{}, errors.New(""))
	_, err := s.service.GetPrice()
	s.Error(err)
}

func (s *Suite) TestMinorEmptyAsks() {
	s.client.On("DoRequest").Return(&models.PriceDepth{
		Asks: nil,
		Bids: nil,
	}, nil)
	_, err := s.service.GetPrice()
	s.Error(err)
	s.Equal(err.Error(), "zero len priceDepth.Asks")
}

func (s *Suite) TestMinorEmptyBids() {
	s.client.On("DoRequest").Return(&models.PriceDepth{
		Asks: make([]models.Ask, 1),
		Bids: nil,
	}, nil)
	_, err := s.service.GetPrice()
	s.Error(err)
	s.Equal(err.Error(), "zero len priceDepth.Bids")
}

func (s *Suite) TestMinorInvalidFirstAsk() {
	s.client.On("DoRequest").Return(&models.PriceDepth{
		Asks: []models.Ask{{Price: "not number value"}},
		Bids: make([]models.Bid, 1),
	}, nil)
	_, err := s.service.GetPrice()
	s.Error(err)
}

func (s *Suite) TestMinorInvalidFirsBid() {
	s.client.On("DoRequest").Return(&models.PriceDepth{
		Asks: []models.Ask{{Price: "67.15"}},
		Bids: []models.Bid{{Price: "not number value"}},
	}, nil)
	_, err := s.service.GetPrice()
	s.Error(err)
}
