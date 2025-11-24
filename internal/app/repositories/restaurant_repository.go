package repositories

import (
	"eatright-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RestaurantRepository interface defines restaurant data access methods
type RestaurantRepository interface {
	Create(restaurant *models.Restaurant) error
	FindByID(id uuid.UUID) (*models.Restaurant, error)
	FindAll() ([]models.Restaurant, error)
	FindByOwnerID(ownerID uuid.UUID) ([]models.Restaurant, error)
	FindNearby(lat, lng, maxDistance float64) ([]models.Restaurant, error)
	Update(restaurant *models.Restaurant) error
}

// restaurantRepository implements RestaurantRepository
type restaurantRepository struct {
	db *gorm.DB
}

// NewRestaurantRepository creates a new restaurant repository
func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{db: db}
}

// Create creates a new restaurant
func (r *restaurantRepository) Create(restaurant *models.Restaurant) error {
	return r.db.Create(restaurant).Error
}

// FindByID finds a restaurant by ID with owner and listings preloaded
func (r *restaurantRepository) FindByID(id uuid.UUID) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := r.db.Preload("Owner").Preload("Listings").Where("id = ?", id).First(&restaurant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &restaurant, nil
}

// FindAll finds all restaurants
func (r *restaurantRepository) FindAll() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.Preload("Owner").Find(&restaurants).Error
	return restaurants, err
}

// FindByOwnerID finds restaurants by owner ID
func (r *restaurantRepository) FindByOwnerID(ownerID uuid.UUID) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.Where("owner_id = ?", ownerID).Find(&restaurants).Error
	return restaurants, err
}

// FindNearby finds restaurants within a certain distance
// Note: This is a simple implementation. For production, consider using PostGIS
func (r *restaurantRepository) FindNearby(lat, lng, maxDistance float64) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	// For now, return all restaurants. The service layer will calculate distances
	// In production, use PostGIS for efficient geospatial queries
	err := r.db.Preload("Owner").Find(&restaurants).Error
	return restaurants, err
}

// Update updates a restaurant
func (r *restaurantRepository) Update(restaurant *models.Restaurant) error {
	return r.db.Save(restaurant).Error
}
