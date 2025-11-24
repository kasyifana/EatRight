package handlers

import (
	"eatright-backend/internal/app/middlewares"
	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/services"
	"eatright-backend/internal/app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// OrderHandler handles order endpoints
type OrderHandler struct {
	orderService services.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	ListingID uuid.UUID `json:"listing_id"`
	Qty       int       `json:"qty"`
}

// CreateOrder creates a new order
// @Summary Create order
// @Description Creates a new food order and automatically decrements stock
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrderRequest true "Order Details"
// @Success 201 {object} utils.Response{data=models.Order} "Order created successfully"
// @Failure 400 {object} utils.Response "Invalid request or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate input
	if req.ListingID == uuid.Nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "listing_id is required", nil)
	}
	if req.Qty <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "qty must be greater than 0", nil)
	}

	// Create order
	order := &models.Order{
		UserID:    userID,
		ListingID: req.ListingID,
		Qty:       req.Qty,
	}

	if err := h.orderService.CreateOrder(order); err != nil {
		if err == models.ErrInsufficientStock {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Insufficient stock available", err)
		}
		if err == models.ErrInvalidInput {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Listing is not active", err)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create order", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Order created successfully", order)
}

// GetMyOrders retrieves all orders for the authenticated user
// @Summary Get user's order history
// @Description Retrieves all orders placed by the authenticated user
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]models.Order} "Orders retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders/me [get]
func (h *OrderHandler) GetMyOrders(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get orders", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Orders retrieved successfully", orders)
}

// UpdateOrderStatusRequest represents the request body for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status"` // "pending", "ready", "completed", "cancelled"
}

// UpdateOrderStatus updates the status of an order
// @Summary Update order status
// @Description Updates order status - pending/ready/completed/cancelled (restaurant owner only)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID (UUID)"
// @Param request body UpdateOrderStatusRequest true "Status Update"
// @Success 200 {object} utils.Response "Order status updated successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 403 {object} utils.Response "Forbidden - not restaurant owner"
// @Router /orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	// Get user ID from context (must be restaurant owner)
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Get order ID from params
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid order ID", err)
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate status
	status := models.OrderStatus(req.Status)
	if status != models.OrderStatusPending &&
		status != models.OrderStatusReady &&
		status != models.OrderStatusCompleted &&
		status != models.OrderStatusCancelled {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid status", nil)
	}

	// Update status
	if err := h.orderService.UpdateOrderStatus(id, status, userID); err != nil {
		if err == models.ErrUnauthorized {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not the owner of this restaurant", err)
		}
		if err == models.ErrInvalidStatusTransition {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid status transition", err)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update order status", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Order status updated successfully", nil)
}
