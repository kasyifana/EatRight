package middlewares

import (
	"eatright-backend/internal/app/config"
	"eatright-backend/internal/app/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT token and attaches user claims to context
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		fmt.Println("DEBUG [Middleware] - Authorization Header:", authHeader)

		if authHeader == "" {
			fmt.Println("DEBUG [Middleware] - Missing Authorization Header")
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authorization header required", nil)
		}

		// Extract token
		tokenString, err := utils.ExtractToken(authHeader)
		if err != nil {
			fmt.Println("DEBUG [Middleware] - Extract Token Error:", err)
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization header", err)
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			fmt.Println("DEBUG [Middleware] - Validate Token Error:", err)
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token", err)
		}

		fmt.Println("DEBUG [Middleware] - Token Valid. UserID:", claims.UserID, "Role:", claims.Role)

		// Store claims in context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}
	return userID, nil
}

// GetUserEmail retrieves user email from context
func GetUserEmail(c *fiber.Ctx) (string, error) {
	email, ok := c.Locals("userEmail").(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}
	return email, nil
}

// GetUserRole retrieves user role from context
func GetUserRole(c *fiber.Ctx) (string, error) {
	role, ok := c.Locals("userRole").(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}
	return role, nil
}
