package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a food order
type Order struct {
	ID         uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	ListingID  uuid.UUID   `gorm:"type:uuid;not null;index" json:"listing_id"`
	Qty        int         `gorm:"not null" json:"qty"`
	TotalPrice int         `gorm:"not null" json:"total_price"` // Total price in smallest currency unit
	Status     OrderStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Listing Listing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}

// BeforeCreate hook to generate UUID and set defaults
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.Status == "" {
		o.Status = OrderStatusPending
	}
	return nil
}

// TableName specifies the table name for Order model
func (Order) TableName() string {
	return "orders"
}

// CanUpdateStatus checks if the order can transition to the new status
func (o *Order) CanUpdateStatus(newStatus OrderStatus) bool {
	// Cancelled orders cannot be updated
	if o.Status == OrderStatusCancelled {
		return false
	}
	// Completed orders cannot be updated
	if o.Status == OrderStatusCompleted {
		return false
	}
	return true
}

// UpdateStatus updates the order status if valid
func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	if !o.CanUpdateStatus(newStatus) {
		return ErrInvalidStatusTransition
	}
	o.Status = newStatus
	return nil
}
