package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	authpb "github.com/my-store/pkg/api/auth"
	orderpb "github.com/my-store/pkg/api/order"
)

// Middleware to validate JWT token
func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Validate token with Auth Service
		resp, err := s.clients.Auth.Validate(ctx, &authpb.ValidateRequest{Token: token})
		if err != nil || resp.Status != 0 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store user ID in context
		ctx = context.WithValue(r.Context(), "userID", resp.UserId)
		next(w, r.WithContext(ctx))
	}
}

func (s *Server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Items []struct {
			ProductID int64   `json:"product_id"`
			Quantity  int32   `json:"quantity"`
			Price     float64 `json:"price"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to gRPC format
	var orderItems []*orderpb.OrderItem
	for _, item := range req.Items {
		orderItems = append(orderItems, &orderpb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.clients.Order.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId: userID,
		Items:  orderItems,
	})

	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"order_id": resp.OrderId})
}

