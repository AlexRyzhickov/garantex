package handler

import (
	"context"
	"garantex/internal/models"
	"garantex/internal/pb"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PriceService interface {
	GetPrice() (models.Price, error)
	SavePrice(price models.Price) error
}

type Handler struct {
	s      PriceService
	logger *zap.Logger
	pb.UnimplementedCryptoExchangeServiceServer
}

func New(s PriceService, logger *zap.Logger) *Handler {
	return &Handler{
		s:      s,
		logger: logger,
	}
}

func (h *Handler) GetPrice(_ context.Context, _ *pb.GetPriceRequest) (*pb.GetPriceResponse, error) {
	price, err := h.s.GetPrice()
	if err != nil {
		h.logger.Error("GetPrice method error", zap.Error(err))
		return nil, status.Error(codes.Internal, "Internal error")
	}

	err = h.s.SavePrice(price)
	if err != nil {
		h.logger.Error("SavePrice method error", zap.Error(err))
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &pb.GetPriceResponse{
		Ts:       price.Timestamp,
		AskPrice: price.AskPrice,
		BidPrice: price.BidPrice,
	}, nil
}
