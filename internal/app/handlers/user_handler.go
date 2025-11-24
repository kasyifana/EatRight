package handlers

import (
	"eatright-backend/internal/app/middlewares"
	"eatright-backend/internal/app/services"
	"eatright-backend/internal/app/utils"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user endpoints
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe retrieves the authenticated user's profile
// GET /api/users/me
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Get user
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User retrieved successfully", user)
}
