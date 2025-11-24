package handlers

import (
	"fmt"
	"strconv"

	"eatright-backend/internal/app/middlewares"
	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/services"
	"eatright-backend/internal/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RestaurantHandler handles restaurant endpoints
type RestaurantHandler struct {
	restaurantService services.RestaurantService
}

// NewRestaurantHandler creates a new restaurant handler
func NewRestaurantHandler(restaurantService services.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
	}
}

// CreateRestaurantRequest represents the request body for creating a restaurant
type CreateRestaurantRequest struct {
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	ClosingTime string  `json:"closing_time"` // Format: "HH:MM:SS"
}

// CreateRestaurant creates a new restaurant
// @Summary Create restaurant
// @Description Creates a new restaurant (requires restaurant role)
// @Tags Restaurants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRestaurantRequest true "Restaurant Details"
// @Success 201 {object} utils.Response{data=models.Restaurant} "Restaurant created successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 403 {object} utils.Response "Forbidden - restaurant role required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /restaurants [post]
func (h *RestaurantHandler) CreateRestaurant(c *fiber.Ctx) error {
	// DEBUG: Check headers and context
	authHeader := c.Get("Authorization")
	fmt.Println("DEBUG - Header Auth:", authHeader)

	userContext := c.Locals("user")
	fmt.Println("DEBUG - User Context:", userContext)

	// Get user ID from context
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	var req CreateRestaurantRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate input
	if req.Name == "" || req.Address == "" || req.ClosingTime == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing required fields", nil)
	}

	// Parse closing time
	closingTime := models.TimeOnly{}
	if err := closingTime.UnmarshalJSON([]byte(`"` + req.ClosingTime + `"`)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid closing_time format (expected HH:MM:SS)", err)
	}

	// Create restaurant
	restaurant := &models.Restaurant{
		Name:        req.Name,
		Address:     req.Address,
		Lat:         req.Lat,
		Lng:         req.Lng,
		ClosingTime: closingTime,
	}

	if err := h.restaurantService.CreateRestaurant(restaurant, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create restaurant", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Restaurant created successfully", restaurant)
}

// GetRestaurants retrieves restaurants, optionally filtered by proximity
// @Summary List restaurants
// @Description Retrieves all restaurants or nearby restaurants if lat/lng provided
// @Tags Restaurants
// @Accept json
// @Produce json
// @Param lat query number false "Latitude for nearby search"
// @Param lng query number false "Longitude for nearby search"
// @Param distance query number false "Max distance in km (default: 10)"
// @Success 200 {object} utils.Response{data=[]models.Restaurant} "Restaurants retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid parameters"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /restaurants [get]
func (h *RestaurantHandler) GetRestaurants(c *fiber.Ctx) error {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	distanceStr := c.Query("distance", "10") // Default 10km

	// If lat/lng provided, get nearby restaurants
	if latStr != "" && lngStr != "" {
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid lat parameter", err)
		}

		lng, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid lng parameter", err)
		}

		distance, err := strconv.ParseFloat(distanceStr, 64)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid distance parameter", err)
		}

		restaurants, err := h.restaurantService.GetNearbyRestaurants(lat, lng, distance)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get nearby restaurants", err)
		}

		return utils.SuccessResponse(c, fiber.StatusOK, "Restaurants retrieved successfully", restaurants)
	}

	// Otherwise get all restaurants
	restaurants, err := h.restaurantService.GetAllRestaurants()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get restaurants", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Restaurants retrieved successfully", restaurants)
}

// GetRestaurantByID retrieves a restaurant by ID
// @Summary Get restaurant by ID
// @Description Retrieves detailed information about a specific restaurant
// @Tags Restaurants
// @Accept json
// @Produce json
// @Param id path string true "Restaurant ID (UUID)"
// @Success 200 {object} utils.Response{data=models.Restaurant} "Restaurant retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid restaurant ID"
// @Failure 404 {object} utils.Response "Restaurant not found"
// @Router /restaurants/{id} [get]
func (h *RestaurantHandler) GetRestaurantByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid restaurant ID", err)
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Restaurant not found", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Restaurant retrieved successfully", restaurant)
}
