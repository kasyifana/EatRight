package services

import (
	"sort"

	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/repositories"
	"eatright-backend/internal/app/utils"

	"github.com/google/uuid"
)

// RestaurantService handles restaurant-related business logic
type RestaurantService interface {
	CreateRestaurant(restaurant *models.Restaurant, ownerID uuid.UUID) error
	GetRestaurantByID(id uuid.UUID) (*models.Restaurant, error)
	GetAllRestaurants() ([]models.Restaurant, error)
	GetNearbyRestaurants(lat, lng, maxDistance float64) ([]models.Restaurant, error)
	UpdateRestaurant(restaurant *models.Restaurant) error
}

// restaurantService implements RestaurantService
type restaurantService struct {
	restaurantRepo repositories.RestaurantRepository
	userRepo       repositories.UserRepository
}

// NewRestaurantService creates a new restaurant service
func NewRestaurantService(restaurantRepo repositories.RestaurantRepository, userRepo repositories.UserRepository) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
		userRepo:       userRepo,
	}
}

// CreateRestaurant creates a new restaurant
func (s *restaurantService) CreateRestaurant(restaurant *models.Restaurant, ownerID uuid.UUID) error {
	// Verify owner exists and has restaurant role
	owner, err := s.userRepo.FindByID(ownerID)
	if err != nil {
		return err
	}

	if !owner.IsRestaurant() {
		return models.ErrUnauthorized
	}

	restaurant.OwnerID = ownerID
	return s.restaurantRepo.Create(restaurant)
}

// GetRestaurantByID retrieves a restaurant by ID
func (s *restaurantService) GetRestaurantByID(id uuid.UUID) (*models.Restaurant, error) {
	return s.restaurantRepo.FindByID(id)
}

// GetAllRestaurants retrieves all restaurants
func (s *restaurantService) GetAllRestaurants() ([]models.Restaurant, error) {
	return s.restaurantRepo.FindAll()
}

// GetNearbyRestaurants retrieves restaurants within maxDistance (in km)
func (s *restaurantService) GetNearbyRestaurants(lat, lng, maxDistance float64) ([]models.Restaurant, error) {
	// Get all restaurants (in production, use PostGIS for better performance)
	restaurants, err := s.restaurantRepo.FindNearby(lat, lng, maxDistance)
	if err != nil {
		return nil, err
	}

	// Calculate distances and filter
	type restaurantWithDistance struct {
		restaurant models.Restaurant
		distance   float64
	}

	var nearby []restaurantWithDistance
	for _, restaurant := range restaurants {
		distance := utils.CalculateDistance(lat, lng, restaurant.Lat, restaurant.Lng)
		if distance <= maxDistance {
			nearby = append(nearby, restaurantWithDistance{
				restaurant: restaurant,
				distance:   distance,
			})
		}
	}

	// Sort by distance
	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].distance < nearby[j].distance
	})

	// Extract restaurants
	result := make([]models.Restaurant, len(nearby))
	for i, item := range nearby {
		result[i] = item.restaurant
	}

	return result, nil
}

// UpdateRestaurant updates a restaurant
func (s *restaurantService) UpdateRestaurant(restaurant *models.Restaurant) error {
	return s.restaurantRepo.Update(restaurant)
}
