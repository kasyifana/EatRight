package utils

import "github.com/gofiber/fiber/v2"

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *fiber.Ctx, errors []string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "Validation failed",
		"errors":  errors,
	})
}
