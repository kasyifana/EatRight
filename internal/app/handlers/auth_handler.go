package handlers

import (
	"eatright-backend/internal/app/services"
	"eatright-backend/internal/app/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// VerifyTokenRequest represents the request body for token verification
type VerifyTokenRequest struct {
	SupabaseToken string `json:"supabase_token"`
}

// VerifyTokenResponse represents the response for token verification
type VerifyTokenResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// VerifyToken verifies a Supabase token and returns a backend JWT
// @Summary Verify Supabase token
// @Description Verifies a Supabase Google OAuth token and returns a backend JWT for API authentication
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body VerifyTokenRequest true "Supabase Token"
// @Success 200 {object} utils.Response{data=VerifyTokenResponse} "Authentication successful"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 401 {object} utils.Response "Token verification failed"
// @Router /auth/verify [post]
func (h *AuthHandler) VerifyToken(c *fiber.Ctx) error {
	var req VerifyTokenRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate input
	if req.SupabaseToken == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "supabase_token is required", nil)
	}

	// Verify token and get/create user
	user, jwtToken, err := h.authService.VerifySupabaseToken(req.SupabaseToken)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token verification failed", err)
	}

	// Return JWT and user info
	response := VerifyTokenResponse{
		Token: jwtToken,
		User:  user,
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Authentication successful", response)
}
