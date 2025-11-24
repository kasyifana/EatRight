package middlewares

import (
	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/utils"

	"github.com/gofiber/fiber/v2"
)

// RestaurantOnly ensures only users with restaurant role can access the endpoint
func RestaurantOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, err := GetUserRole(c)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
		}

		if role != string(models.RoleRestaurant) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Restaurant role required", nil)
		}

		return c.Next()
	}
}
