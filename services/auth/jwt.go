package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define a secret key for signing tokens.
// In a real application, this should come from environment variables.
var jwtSecret = []byte("super-secret-key")

// Claims defines the structure of the JWT payload.
type Claims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a given user ID.
func GenerateToken(userId int64) (string, error) {
	// Create the claims
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token with HS256 signing algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	return token.SignedString(jwtSecret)
}

// ValidateToken parses and validates a JWT token.
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		// It's crucial to verify the alg is what you expect to prevent downgrade attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

