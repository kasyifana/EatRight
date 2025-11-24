package handlers

import (
	"eatright-backend/internal/app/middlewares"
	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/services"
	"eatright-backend/internal/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ListingHandler handles listing endpoints
type ListingHandler struct {
	listingService services.ListingService
}

// NewListingHandler creates a new listing handler
func NewListingHandler(listingService services.ListingService) *ListingHandler {
	return &ListingHandler{
		listingService: listingService,
	}
}

// CreateListingRequest represents the request body for creating a listing
type CreateListingRequest struct {
	Type        string  `json:"type"` // "mystery_box" or "reveal"
	Name        *string `json:"name"` // Optional for mystery box
	Description string  `json:"description"`
	Price       int     `json:"price"`
	Stock       int     `json:"stock"`
	PhotoURL    string  `json:"photo_url"`
	PickupTime  string  `json:"pickup_time"` // Format: "HH:MM:SS"
}

// CreateListing creates a new listing for a restaurant
// @Summary Create food listing
// @Description Creates a new food listing for a restaurant (mystery box or reveal item)
// @Tags Listings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Restaurant ID (UUID)"
// @Param request body CreateListingRequest true "Listing Details"
// @Success 201 {object} utils.Response{data=models.Listing} "Listing created successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Router /restaurants/{id}/listings [post]
func (h *ListingHandler) CreateListing(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Get restaurant ID from params
	restaurantIDStr := c.Params("id")
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid restaurant ID", err)
	}

	var req CreateListingRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate input
	if req.Type == "" || req.Description == "" || req.PickupTime == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing required fields", nil)
	}

	// Validate type
	listingType := models.ListingType(req.Type)
	if listingType != models.ListingTypeMysteryBox && listingType != models.ListingTypeReveal {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing type (must be 'mystery_box' or 'reveal')", nil)
	}

	// Parse pickup time
	pickupTime := models.TimeOnly{}
	if err := pickupTime.UnmarshalJSON([]byte(`"` + req.PickupTime + `"`)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid pickup_time format (expected HH:MM:SS)", err)
	}

	// Create listing
	listing := &models.Listing{
		RestaurantID: restaurantID,
		Type:         listingType,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		Stock:        req.Stock,
		PhotoURL:     req.PhotoURL,
		PickupTime:   pickupTime,
		IsActive:     true,
	}

	if err := h.listingService.CreateListing(listing, userID); err != nil {
		if err == models.ErrUnauthorized {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not the owner of this restaurant", err)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create listing", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Listing created successfully", listing)
}

// GetListings retrieves all active listings
// @Summary List all food offerings
// @Description Retrieves all active food listings from all restaurants
// @Tags Listings
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Listing} "Listings retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /listings [get]
func (h *ListingHandler) GetListings(c *fiber.Ctx) error {
	listings, err := h.listingService.GetAllListings(true) // Only active listings
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get listings", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Listings retrieved successfully", listings)
}

// GetListingByID retrieves a listing by ID
// @Summary Get listing details
// @Description Retrieves detailed information about a specific food listing
// @Tags Listings
// @Accept json
// @Produce json
// @Param id path string true "Listing ID (UUID)"
// @Success 200 {object} utils.Response{data=models.Listing} "Listing retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid listing ID"
// @Failure 404 {object} utils.Response "Listing not found"
// @Router /listings/{id} [get]
func (h *ListingHandler) GetListingByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID", err)
	}

	listing, err := h.listingService.GetListingByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Listing not found", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Listing retrieved successfully", listing)
}

// UpdateStockRequest represents the request body for updating stock
type UpdateStockRequest struct {
	Quantity int `json:"quantity"` // Can be positive (add) or negative (reduce)
}

// UpdateStock updates the stock of a listing
// @Summary Update listing stock
// @Description Updates the stock quantity of a listing (restaurant owner only)
// @Tags Listings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Listing ID (UUID)"
// @Param request body UpdateStockRequest true "Stock Update"
// @Success 200 {object} utils.Response "Stock updated successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Router /listings/{id}/stock [patch]
func (h *ListingHandler) UpdateStock(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Get listing ID from params
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID", err)
	}

	var req UpdateStockRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Update stock
	if err := h.listingService.UpdateStock(id, req.Quantity, userID); err != nil {
		if err == models.ErrUnauthorized {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not the owner of this listing", err)
		}
		if err == models.ErrNegativeStock {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Stock cannot be negative", err)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update stock", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Stock updated successfully", nil)
}

// UpdateStatusRequest represents the request body for updating listing status
type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// UpdateStatus toggles the active status of a listing
// @Summary Toggle listing status
// @Description Activates or deactivates a listing (restaurant owner only)
// @Tags Listings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Listing ID (UUID)"
// @Param request body UpdateStatusRequest true "Status Update"
// @Success 200 {object} utils.Response "Status updated successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 403 {object} utils.Response "Forbidden"
// @Router /listings/{id}/status [patch]
func (h *ListingHandler) UpdateStatus(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Get listing ID from params
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID", err)
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Update status
	if err := h.listingService.ToggleActive(id, req.IsActive, userID); err != nil {
		if err == models.ErrUnauthorized {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not the owner of this listing", err)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update status", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Status updated successfully", nil)
}
