package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Auth Service listening on port %d", port)
	
	// gRPC server setup will go here once protos are generated
	// s := grpc.NewServer()
	// pb.RegisterAuthServiceServer(s, &server{})
	
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
	
	// Keep alive for now
	select {}
}

