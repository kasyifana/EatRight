package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleRestaurant UserRole = "restaurant"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Role      UserRole  `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Restaurants []Restaurant `gorm:"foreignKey:OwnerID" json:"restaurants,omitempty"`
	Orders      []Order      `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// BeforeCreate hook to generate UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Role == "" {
		u.Role = RoleUser
	}
	return nil
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// IsRestaurant checks if the user has restaurant role
func (u *User) IsRestaurant() bool {
	return u.Role == RoleRestaurant
}
