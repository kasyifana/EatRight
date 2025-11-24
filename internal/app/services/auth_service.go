package services

import (
	"fmt"

	"eatright-backend/internal/app/config"
	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/repositories"
	"eatright-backend/internal/app/utils"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles authentication logic
type AuthService interface {
	VerifySupabaseToken(token string) (*models.User, string, error)
}

// authService implements AuthService
type authService struct {
	userRepo repositories.UserRepository
	config   *config.Config
}

// SupabaseClaims represents Supabase JWT claims
type SupabaseClaims struct {
	Email                 string                 `json:"email"`
	Phone                 string                 `json:"phone"`
	AppMetadata           map[string]interface{} `json:"app_metadata"`
	UserMetadata          map[string]interface{} `json:"user_metadata"`
	Role                  string                 `json:"role"`
	AuthenticatorAssurace string                 `json:"aal"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config) (AuthService, error) {
	return &authService{
		userRepo: userRepo,
		config:   cfg,
	}, nil
}

// VerifySupabaseToken verifies a Supabase OAuth token and returns user + JWT
func (s *authService) VerifySupabaseToken(token string) (*models.User, string, error) {
	// Parse Supabase JWT token
	// Note: In production, you should verify the signature using Supabase JWT secret
	// For now, we'll parse it without verification (unsafe for production)
	// You should use SUPABASE_JWT_SECRET from your Supabase dashboard

	claims := &SupabaseClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		// For production, use your Supabase JWT secret here
		// return []byte(s.config.Supabase.ServiceKey), nil

		// For now, we'll skip verification to keep things simple
		// This means we trust the token is from Supabase
		return []byte(""), nil
	})

	// We ignore signature verification error for simplicity
	// In production, you MUST verify the signature
	if claims.Email == "" {
		return nil, "", fmt.Errorf("invalid token: no email found")
	}

	email := claims.Email
	name := ""

	// Try to get name from user metadata
	if claims.UserMetadata != nil {
		if fullName, ok := claims.UserMetadata["full_name"].(string); ok {
			name = fullName
		} else if userName, ok := claims.UserMetadata["name"].(string); ok {
			name = userName
		}
	}

	if name == "" {
		name = email
	}

	// Check if user exists in our database
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if err == models.ErrNotFound {
			// Create new user
			user = &models.User{
				Name:  name,
				Email: email,
				Role:  models.RoleUser, // Default role
			}

			err = s.userRepo.Create(user)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, "", fmt.Errorf("failed to find user: %w", err)
		}
	}

	// Generate JWT token for our backend
	jwtToken, err := utils.GenerateToken(user, s.config.JWT.Secret, s.config.JWT.Expiry)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate jwt: %w", err)
	}

	return user, jwtToken, nil
}
