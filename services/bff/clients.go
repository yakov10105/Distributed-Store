package main

import (
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/my-store/pkg/api/auth"
	orderpb "github.com/my-store/pkg/api/order"
)

type ServiceClients struct {
	Auth  authpb.AuthServiceClient
	Order orderpb.OrderServiceClient
}

func InitClients() (*ServiceClients, error) {
	// Connect to Auth Service
	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}
	authConn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to Auth Service at %s", authAddr)

	// Connect to Order Service
	orderAddr := os.Getenv("ORDER_SERVICE_ADDR")
	if orderAddr == "" {
		orderAddr = "localhost:50052"
	}
	orderConn, err := grpc.NewClient(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to Order Service at %s", orderAddr)

	return &ServiceClients{
		Auth:  authpb.NewAuthServiceClient(authConn),
		Order: orderpb.NewOrderServiceClient(orderConn),
	}, nil
}

