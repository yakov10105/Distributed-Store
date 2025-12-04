package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := 50052
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Order Service listening on port %d", port)
	
	// Keep alive for now
	select {}
}

