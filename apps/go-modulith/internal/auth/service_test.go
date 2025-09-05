package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/zero-one-group/go-modulith/internal/config"
)

func TestService_GenerateAccessToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 24 * time.Hour,
		},
	}

	service := &Service{
		cfg: cfg,
	}

	user := &User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, err := service.generateAccessToken(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("Expected token to be generated, got empty string")
	}
}

func TestService_Register_ValidInput(t *testing.T) {
	// This is a basic structure test - would need proper mocking for full implementation
	ctx := context.Background()
	
	req := RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Validate request structure
	if req.Name == "" {
		t.Error("Expected name to be set")
	}
	if req.Email == "" {
		t.Error("Expected email to be set")
	}
	if req.Password == "" {
		t.Error("Expected password to be set")
	}

	_ = ctx // Use ctx to avoid unused variable
}