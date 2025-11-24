package repositories

import (
	"eatright-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderRepository interface defines order data access methods
type OrderRepository interface {
	Create(order *models.Order, listing *models.Listing) error
	FindByID(id uuid.UUID) (*models.Order, error)
	FindByUserID(userID uuid.UUID) ([]models.Order, error)
	FindByRestaurantID(restaurantID uuid.UUID) ([]models.Order, error)
	UpdateStatus(id uuid.UUID, status models.OrderStatus) error
}

// orderRepository implements OrderRepository
type orderRepository struct {
	db          *gorm.DB
	listingRepo ListingRepository
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB, listingRepo ListingRepository) OrderRepository {
	return &orderRepository{
		db:          db,
		listingRepo: listingRepo,
	}
}

// Create creates a new order and decrements stock in a transaction
func (r *orderRepository) Create(order *models.Order, listing *models.Listing) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Validate quantity
		if order.Qty <= 0 {
			return models.ErrInvalidQuantity
		}

		// Check stock availability and decrement (with row lock)
		err := r.listingRepo.UpdateStockWithTx(tx, order.ListingID, -order.Qty)
		if err != nil {
			if err == models.ErrNegativeStock {
				return models.ErrInsufficientStock
			}
			return err
		}

		// Calculate total price
		order.TotalPrice = listing.Price * order.Qty

		// Create the order
		err = tx.Create(order).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// FindByID finds an order by ID with related data preloaded
func (r *orderRepository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("Listing").Preload("Listing.Restaurant").Where("id = ?", id).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

// FindByUserID finds all orders by user ID
func (r *orderRepository) FindByUserID(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Listing").Preload("Listing.Restaurant").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// FindByRestaurantID finds all orders for a restaurant
func (r *orderRepository) FindByRestaurantID(restaurantID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("User").Preload("Listing").
		Joins("JOIN listings ON listings.id = orders.listing_id").
		Where("listings.restaurant_id = ?", restaurantID).
		Order("orders.created_at DESC").
		Find(&orders).Error
	return orders, err
}

// UpdateStatus updates the status of an order
func (r *orderRepository) UpdateStatus(id uuid.UUID, status models.OrderStatus) error {
	var order models.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.ErrNotFound
		}
		return err
	}

	// Validate status transition
	if err := order.UpdateStatus(status); err != nil {
		return err
	}

	return r.db.Model(&order).Update("status", status).Error
}
