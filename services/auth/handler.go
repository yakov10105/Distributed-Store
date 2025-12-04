package main

import (
	"context"

	pb "github.com/my-store/pkg/api/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer implements the generated AuthServiceServer interface.
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	store *UserStore
}

// NewAuthServer creates a new instance of our gRPC server.
func NewAuthServer(store *UserStore) *AuthServer {
	return &AuthServer{
		store: store,
	}
}

// Register handles user registration.
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &pb.RegisterResponse{
			Status: int32(codes.InvalidArgument),
			Error:  "Email and password are required",
		}, nil
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	_, err = s.store.Create(req.Email, string(hashedBytes))
	if err != nil {
		// In a real app, check for specific DB error code (e.g., unique constraint violation)
		// For now, we assume most errors are duplicates or connection issues
		return &pb.RegisterResponse{
			Status: int32(codes.AlreadyExists),
			Error:  "User already exists or database error",
		}, nil
	}

	return &pb.RegisterResponse{
		Status: int32(codes.OK),
	}, nil
}

// Login handles user authentication.
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.store.FindByEmail(req.Email)
	if err != nil {
		return &pb.LoginResponse{
			Status: int32(codes.NotFound),
			Error:  "User not found",
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &pb.LoginResponse{
			Status: int32(codes.Unauthenticated),
			Error:  "Invalid credentials",
		}, nil
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}

	return &pb.LoginResponse{
		Status: int32(codes.OK),
		Token:  token,
	}, nil
}

// Validate checks if a token is valid.
func (s *AuthServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Status: int32(codes.Unauthenticated),
			Error:  "Invalid token",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: int32(codes.OK),
		UserId: claims.UserId,
	}, nil
}
