package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// User represents the authenticated user context
type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

// contextKey is used for storing user in context
type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware creates a middleware for Supabase JWT authentication
func AuthMiddleware(supabaseJWTSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// For now, we'll create a simple JWT verification stub
		// In production, you would use the Supabase GoTrue client
		user, err := verifyJWT(tokenString, supabaseJWTSecret)
		if err != nil {
			log.Error().Err(err).Msg("Failed to verify JWT")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store the user in the context
		ctx := context.WithValue(c.Request.Context(), userContextKey, user)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// verifyJWT is a simplified JWT verification (in production, use proper JWT library)
func verifyJWT(tokenString, _secret string) (*User, error) {
	// This is a simplified implementation
	// In production, you should use a proper JWT library like jwt-go
	// and integrate with Supabase's GoTrue service

	// For development, we'll extract user info from a mock token
	// In production, this would verify the JWT signature and extract claims

	// Mock implementation - extract user ID from token (for development)
	// In production, replace with proper JWT verification
	if tokenString == "" {
		return nil, fmt.Errorf("empty token")
	}

	// Create a mock user for development
	// In production, this would come from verified JWT claims
	var user *User

	// Use a deterministic user ID for test-token to make testing easier
	if tokenString == "test-token" {
		// Use a fixed UUID for test-token
		testUUID, _ := uuid.Parse("12345678-1234-1234-1234-123456789012")
		user = &User{
			ID:    testUUID, // Fixed UUID for test-token
			Email: "test@example.com", // This should come from JWT claims
			Name:  "Test User", // This should come from JWT claims
		}
	} else {
		// For other tokens, generate a new UUID (current behavior)
		user = &User{
			ID:    uuid.New(), // This should come from JWT claims
			Email: "user@example.com", // This should come from JWT claims
			Name:  "Test User", // This should come from JWT claims
		}
	}

	log.Debug().Str("user_id", user.ID.String()).Msg("User authenticated")
	return user, nil
}

// GetUser retrieves the authenticated user from the context
func GetUser(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

// GetUserFromGin retrieves the authenticated user from Gin context
func GetUserFromGin(c *gin.Context) (*User, error) {
	return GetUser(c.Request.Context())
}
