package main

import (
	"context"

	pb "github.com/my-store/pkg/api/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderServer implements the generated OrderServiceServer interface.
type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	store *OrderStore
}

// NewOrderServer creates a new instance of our gRPC server.
func NewOrderServer(store *OrderStore) *OrderServer {
	return &OrderServer{
		store: store,
	}
}

// CreateOrder handles order creation.
func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	if len(req.Items) == 0 {
		return &pb.CreateOrderResponse{
			Status: int32(codes.InvalidArgument),
			Error:  "Order must have at least one item",
		}, nil
	}

	order, err := s.store.Create(req.UserId, req.Items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create order: %v", err)
	}

	return &pb.CreateOrderResponse{
		Status:  int32(codes.OK),
		OrderId: order.ID,
	}, nil
}

// GetOrder retrieves order details.
func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.store.Get(req.OrderId)
	if err != nil {
		return &pb.GetOrderResponse{
			Status: int32(codes.NotFound),
			Error:  "Order not found",
		}, nil
	}

	return &pb.GetOrderResponse{
		Status:  int32(codes.OK),
		OrderId: order.ID,
		UserId:  order.UserID,
		Items:   order.Items,
	}, nil
}
