package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

type Server struct {
	clients *ServiceClients
}

func main() {
	// 1. Initialize gRPC Clients
	clients, err := InitClients()
	if err != nil {
		log.Fatalf("Failed to initialize clients: %v", err)
	}

	server := &Server{clients: clients}

	// 2. Setup Router
	mux := http.NewServeMux()

	// Public Endpoints
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("BFF is running"))
	})

	// Auth Endpoints
	mux.HandleFunc("/api/auth/register", server.handleRegister)
	mux.HandleFunc("/api/auth/login", server.handleLogin)

	// Protected Endpoints
	mux.HandleFunc("/api/orders", server.withAuth(server.handleCreateOrder))

	// 3. Setup CORS
	// Allow requests from frontend (localhost:3000)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	// 4. Start Server
	port := "8080"
	log.Printf("BFF Service listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
