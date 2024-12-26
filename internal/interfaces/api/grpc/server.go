package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "marketdata/api/proto"
	"marketdata/internal/application/port/input"
)

type MarketDataServer struct {
	pb.UnimplementedMarketDataServiceServer
	marketDataUseCase input.MarketDataUseCase
}

func NewMarketDataServer(useCase input.MarketDataUseCase) *MarketDataServer {
	return &MarketDataServer{
		marketDataUseCase: useCase,
	}
}

func (s *MarketDataServer) GetOrderBook(ctx context.Context, req *pb.GetOrderBookRequest) (*pb.OrderBook, error) {
	if req.ExchangeId == "" || req.Symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "exchange_id and symbol are required")
	}

	orderbook, err := s.marketDataUseCase.GetOrderBook(ctx, req.ExchangeId, req.Symbol)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if orderbook == nil {
		return nil, status.Error(codes.NotFound, "orderbook not found")
	}

	// Convert DTO to proto message
	return convertToProto(orderbook), nil
}

func (s *MarketDataServer) SubscribeOrderBook(req *pb.SubscribeOrderBookRequest, stream pb.MarketDataService_SubscribeOrderBookServer) error {
	if req.ExchangeId == "" || req.Symbol == "" {
		return status.Error(codes.InvalidArgument, "exchange_id and symbol are required")
	}

	updates, err := s.marketDataUseCase.SubscribeOrderBook(stream.Context(), req.ExchangeId, req.Symbol)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case update := <-updates:
			if err := stream.Send(convertToProto(update)); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}
}
