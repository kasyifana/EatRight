package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListingType represents the type of a food listing
type ListingType string

const (
	ListingTypeMysteryBox ListingType = "mystery_box"
	ListingTypeReveal     ListingType = "reveal"
)

// Listing represents a food listing (mystery box or reveal item)
type Listing struct {
	ID           uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RestaurantID uuid.UUID   `gorm:"type:uuid;not null;index" json:"restaurant_id"`
	Type         ListingType `gorm:"type:varchar(20);not null" json:"type"`
	Name         *string     `gorm:"type:varchar(255)" json:"name"` // Nullable for mystery box
	Description  string      `gorm:"type:text;not null" json:"description"`
	Price        int         `gorm:"not null" json:"price"` // Price in smallest currency unit (e.g., cents)
	Stock        int         `gorm:"not null;default:0" json:"stock"`
	PhotoURL     string      `gorm:"type:text" json:"photo_url"`
	PickupTime   TimeOnly    `gorm:"type:time;not null" json:"pickup_time"`
	IsActive     bool        `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	Orders     []Order    `gorm:"foreignKey:ListingID" json:"orders,omitempty"`
}

// BeforeCreate hook to generate UUID before creating
func (l *Listing) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Listing model
func (Listing) TableName() string {
	return "listings"
}

// IsMysteryBox checks if the listing is a mystery box
func (l *Listing) IsMysteryBox() bool {
	return l.Type == ListingTypeMysteryBox
}

// HasStock checks if the listing has available stock
func (l *Listing) HasStock(qty int) bool {
	return l.Stock >= qty
}

// DecrementStock reduces the stock by the given quantity
func (l *Listing) DecrementStock(qty int) error {
	if !l.HasStock(qty) {
		return ErrInsufficientStock
	}
	l.Stock -= qty
	return nil
}
