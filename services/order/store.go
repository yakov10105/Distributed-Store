package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	pb "github.com/my-store/pkg/api/order"
)

// Order represents an order in our system.
type Order struct {
	ID     int64
	UserID int64
	Items  []*pb.OrderItem
	Status string
}

// OrderStore handles database interactions for orders.
type OrderStore struct {
	db *sql.DB
}

// NewOrderStore initializes the store with a database connection.
func NewOrderStore(db *sql.DB) *OrderStore {
	return &OrderStore{db: db}
}

// InitSchema creates the orders table if it doesn't exist.
func (s *OrderStore) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		status TEXT NOT NULL,
		items JSONB NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

// Create adds a new order to the database.
func (s *OrderStore) Create(userID int64, items []*pb.OrderItem) (*Order, error) {
	// Convert items to JSON for storage
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal items: %w", err)
	}

	query := `
		INSERT INTO orders (user_id, status, items) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	var id int64
	err = s.db.QueryRow(query, userID, "PENDING", itemsJSON).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	return &Order{
		ID:     id,
		UserID: userID,
		Items:  items,
		Status: "PENDING",
	}, nil
}

// Get retrieves an order by ID.
func (s *OrderStore) Get(orderID int64) (*Order, error) {
	query := `SELECT id, user_id, status, items FROM orders WHERE id = $1`

	var order Order
	var itemsJSON []byte

	err := s.db.QueryRow(query, orderID).Scan(&order.ID, &order.UserID, &order.Status, &itemsJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return &order, nil
}
