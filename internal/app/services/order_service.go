package services

import (
	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/repositories"

	"github.com/google/uuid"
)

// OrderService handles order-related business logic
type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	GetUserOrders(userID uuid.UUID) ([]models.Order, error)
	GetRestaurantOrders(restaurantID uuid.UUID) ([]models.Order, error)
	UpdateOrderStatus(id uuid.UUID, status models.OrderStatus, requesterID uuid.UUID) error
}

// orderService implements OrderService
type orderService struct {
	orderRepo      repositories.OrderRepository
	listingRepo    repositories.ListingRepository
	restaurantRepo repositories.RestaurantRepository
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo repositories.OrderRepository,
	listingRepo repositories.ListingRepository,
	restaurantRepo repositories.RestaurantRepository,
) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		listingRepo:    listingRepo,
		restaurantRepo: restaurantRepo,
	}
}

// CreateOrder creates a new order
func (s *orderService) CreateOrder(order *models.Order) error {
	// Get listing to check stock and get price
	listing, err := s.listingRepo.FindByID(order.ListingID)
	if err != nil {
		return err
	}

	// Validate listing is active
	if !listing.IsActive {
		return models.ErrInvalidInput
	}

	// Validate quantity
	if order.Qty <= 0 {
		return models.ErrInvalidQuantity
	}

	// Check stock availability
	if !listing.HasStock(order.Qty) {
		return models.ErrInsufficientStock
	}

	// Create order (repository will handle stock decrement in transaction)
	return s.orderRepo.Create(order, listing)
}

// GetOrderByID retrieves an order by ID
func (s *orderService) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	return s.orderRepo.FindByID(id)
}

// GetUserOrders retrieves all orders for a user
func (s *orderService) GetUserOrders(userID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

// GetRestaurantOrders retrieves all orders for a restaurant
func (s *orderService) GetRestaurantOrders(restaurantID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.FindByRestaurantID(restaurantID)
}

// UpdateOrderStatus updates the status of an order
func (s *orderService) UpdateOrderStatus(id uuid.UUID, status models.OrderStatus, requesterID uuid.UUID) error {
	// Get order
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Get restaurant for listing
	restaurant, err := s.restaurantRepo.FindByID(order.Listing.RestaurantID)
	if err != nil {
		return err
	}

	// Verify requester is the restaurant owner
	if restaurant.OwnerID != requesterID {
		return models.ErrUnauthorized
	}

	// Update status
	return s.orderRepo.UpdateStatus(id, status)
}
