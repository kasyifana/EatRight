package services

import (
	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/repositories"

	"github.com/google/uuid"
)

// ListingService handles listing-related business logic
type ListingService interface {
	CreateListing(listing *models.Listing, ownerID uuid.UUID) error
	GetListingByID(id uuid.UUID) (*models.Listing, error)
	GetAllListings(activeOnly bool) ([]models.Listing, error)
	GetListingsByRestaurant(restaurantID uuid.UUID) ([]models.Listing, error)
	UpdateStock(id uuid.UUID, qty int, ownerID uuid.UUID) error
	ToggleActive(id uuid.UUID, active bool, ownerID uuid.UUID) error
}

// listingService implements ListingService
type listingService struct {
	listingRepo    repositories.ListingRepository
	restaurantRepo repositories.RestaurantRepository
}

// NewListingService creates a new listing service
func NewListingService(listingRepo repositories.ListingRepository, restaurantRepo repositories.RestaurantRepository) ListingService {
	return &listingService{
		listingRepo:    listingRepo,
		restaurantRepo: restaurantRepo,
	}
}

// CreateListing creates a new listing
func (s *listingService) CreateListing(listing *models.Listing, ownerID uuid.UUID) error {
	// Verify restaurant exists and belongs to owner
	restaurant, err := s.restaurantRepo.FindByID(listing.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.OwnerID != ownerID {
		return models.ErrUnauthorized
	}

	// Validate stock
	if listing.Stock < 0 {
		return models.ErrNegativeStock
	}

	return s.listingRepo.Create(listing)
}

// GetListingByID retrieves a listing by ID
func (s *listingService) GetListingByID(id uuid.UUID) (*models.Listing, error) {
	return s.listingRepo.FindByID(id)
}

// GetAllListings retrieves all listings
func (s *listingService) GetAllListings(activeOnly bool) ([]models.Listing, error) {
	return s.listingRepo.FindAll(activeOnly)
}

// GetListingsByRestaurant retrieves all listings for a restaurant
func (s *listingService) GetListingsByRestaurant(restaurantID uuid.UUID) ([]models.Listing, error) {
	return s.listingRepo.FindByRestaurantID(restaurantID)
}

// UpdateStock updates the stock of a listing
func (s *listingService) UpdateStock(id uuid.UUID, qty int, ownerID uuid.UUID) error {
	// Get listing with restaurant
	listing, err := s.listingRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verify ownership
	if listing.Restaurant.OwnerID != ownerID {
		return models.ErrUnauthorized
	}

	// Validate stock won't go negative
	if listing.Stock+qty < 0 {
		return models.ErrNegativeStock
	}

	return s.listingRepo.UpdateStock(id, qty)
}

// ToggleActive toggles the active status of a listing
func (s *listingService) ToggleActive(id uuid.UUID, active bool, ownerID uuid.UUID) error {
	// Get listing with restaurant
	listing, err := s.listingRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verify ownership
	if listing.Restaurant.OwnerID != ownerID {
		return models.ErrUnauthorized
	}

	return s.listingRepo.ToggleActive(id, active)
}
