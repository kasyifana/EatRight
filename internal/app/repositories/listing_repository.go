package repositories

import (
	"fmt"

	"eatright-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ListingRepository interface defines listing data access methods
type ListingRepository interface {
	Create(listing *models.Listing) error
	FindByID(id uuid.UUID) (*models.Listing, error)
	FindAll(activeOnly bool) ([]models.Listing, error)
	FindByRestaurantID(restaurantID uuid.UUID) ([]models.Listing, error)
	Update(listing *models.Listing) error
	UpdateStock(id uuid.UUID, qty int) error
	UpdateStockWithTx(tx *gorm.DB, id uuid.UUID, qty int) error
	ToggleActive(id uuid.UUID, active bool) error
}

// listingRepository implements ListingRepository
type listingRepository struct {
	db *gorm.DB
}

// NewListingRepository creates a new listing repository
func NewListingRepository(db *gorm.DB) ListingRepository {
	return &listingRepository{db: db}
}

// Create creates a new listing
func (r *listingRepository) Create(listing *models.Listing) error {
	return r.db.Create(listing).Error
}

// FindByID finds a listing by ID with restaurant preloaded
func (r *listingRepository) FindByID(id uuid.UUID) (*models.Listing, error) {
	var listing models.Listing
	err := r.db.Preload("Restaurant").Preload("Restaurant.Owner").Where("id = ?", id).First(&listing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &listing, nil
}

// FindAll finds all listings, optionally filtering by active status
func (r *listingRepository) FindAll(activeOnly bool) ([]models.Listing, error) {
	var listings []models.Listing
	query := r.db.Preload("Restaurant").Preload("Restaurant.Owner")

	if activeOnly {
		query = query.Where("is_active = ? AND stock > 0", true)
	}

	err := query.Order("created_at DESC").Find(&listings).Error
	return listings, err
}

// FindByRestaurantID finds listings by restaurant ID
func (r *listingRepository) FindByRestaurantID(restaurantID uuid.UUID) ([]models.Listing, error) {
	var listings []models.Listing
	err := r.db.Where("restaurant_id = ?", restaurantID).Order("created_at DESC").Find(&listings).Error
	return listings, err
}

// Update updates a listing
func (r *listingRepository) Update(listing *models.Listing) error {
	return r.db.Save(listing).Error
}

// UpdateStock updates listing stock with validation
func (r *listingRepository) UpdateStock(id uuid.UUID, qty int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.UpdateStockWithTx(tx, id, qty)
	})
}

// UpdateStockWithTx updates listing stock within a transaction
func (r *listingRepository) UpdateStockWithTx(tx *gorm.DB, id uuid.UUID, qty int) error {
	var listing models.Listing

	// Lock the row for update to prevent race conditions
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&listing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.ErrNotFound
		}
		return err
	}

	// Calculate new stock
	newStock := listing.Stock + qty

	// Prevent negative stock
	if newStock < 0 {
		return models.ErrNegativeStock
	}

	// Update stock
	err = tx.Model(&listing).Update("stock", newStock).Error
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

// ToggleActive toggles the active status of a listing
func (r *listingRepository) ToggleActive(id uuid.UUID, active bool) error {
	var listing models.Listing
	err := r.db.Where("id = ?", id).First(&listing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.ErrNotFound
		}
		return err
	}

	return r.db.Model(&listing).Update("is_active", active).Error
}
