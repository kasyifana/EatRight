package services

import (
	"eatright-backend/internal/app/models"
	"eatright-backend/internal/app/repositories"

	"github.com/google/uuid"
)

// UserService handles user-related business logic
type UserService interface {
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

// userService implements UserService
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// UpdateUser updates a user
func (s *userService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}
