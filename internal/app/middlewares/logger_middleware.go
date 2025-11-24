package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger middleware logs HTTP requests
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request details
		log.Printf(
			"%s %s - %d - %v",
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}
