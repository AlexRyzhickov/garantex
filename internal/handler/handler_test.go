package handler

import (
	"context"
	"errors"
	"garantex/internal/mock"
	"garantex/internal/models"
	"garantex/internal/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestService(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	priceService mock.PriceService
	handler      Handler
}

func (s *Suite) SetupTest() {
	s.priceService = mock.PriceService{}
	s.handler = Handler{
		s: &s.priceService,
	}
}

func (s *Suite) TestMajor() {
	s.priceService.On("GetPrice").Return(models.Price{}, nil)
	s.priceService.On("SavePrice", models.Price{}).Return(nil)
	_, err := s.handler.GetPrice(context.Background(), &pb.GetPriceRequest{})
	s.NoError(err)
}

func (s *Suite) TestMinorGetPriceReturnErr() {
	s.priceService.On("GetPrice").Return(models.Price{}, errors.New(""))
	_, err := s.handler.GetPrice(context.Background(), &pb.GetPriceRequest{})
	s.Error(err)
	assert.Equal(s.T(), err.Error(), status.Error(codes.Internal, "Internal error").Error())
}

func (s *Suite) TestMinorSavePriceReturnErr() {
	s.priceService.On("GetPrice").Return(models.Price{}, nil)
	s.priceService.On("SavePrice", models.Price{}).Return(errors.New(""))
	_, err := s.handler.GetPrice(context.Background(), &pb.GetPriceRequest{})
	s.Error(err)
	assert.Equal(s.T(), err.Error(), status.Error(codes.Internal, "Internal error").Error())
}
