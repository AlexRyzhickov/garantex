package handler

import (
	"context"
	"garantex/internal/models"
	"garantex/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type PriceService interface {
	GetPrice() (models.Price, error)
	SavePrice(price models.Price) error
}

type Handler struct {
	s PriceService
	pb.UnimplementedCryptoExchangeServiceServer
}

func New(s PriceService) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) GetPrice(_ context.Context, _ *pb.GetPriceRequest) (*pb.GetPriceResponse, error) {
	price, err := h.s.GetPrice()
	if err != nil {
		log.Println("GetPrice method error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	err = h.s.SavePrice(price)
	if err != nil {
		log.Println("SavePrice method error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &pb.GetPriceResponse{
		Ts:       price.Timestamp,
		AskPrice: price.AskPrice,
		BidPrice: price.BidPrice,
	}, nil
}
