package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx driver for database/sql
)

// User represents a user in our system.
type User struct {
	ID       int64
	Email    string
	Password string // Hashed password
}

// UserStore handles database interactions for users.
type UserStore struct {
	db *sql.DB
}

// NewUserStore initializes the store with a database connection.
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

// InitSchema creates the users table if it doesn't exist.
func (s *UserStore) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

// Create adds a new user to the database.
func (s *UserStore) Create(email, hashedPassword string) (*User, error) {
	query := `
		INSERT INTO users (email, password) 
		VALUES ($1, $2) 
		RETURNING id`

	var id int64
	err := s.db.QueryRow(query, email, hashedPassword).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
	}, nil
}

// FindByEmail retrieves a user by email from the database.
func (s *UserStore) FindByEmail(email string) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`

	var user User
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
