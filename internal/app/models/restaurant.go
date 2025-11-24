package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Restaurant represents a restaurant in the system
type Restaurant struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Address     string    `gorm:"type:text;not null" json:"address"`
	Lat         float64   `gorm:"type:decimal(10,8);not null" json:"lat"`
	Lng         float64   `gorm:"type:decimal(11,8);not null" json:"lng"`
	ClosingTime TimeOnly  `gorm:"type:time;not null" json:"closing_time"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Owner    User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Listings []Listing `gorm:"foreignKey:RestaurantID" json:"listings,omitempty"`
}

// TimeOnly is a custom type for time without date
type TimeOnly struct {
	time.Time
}

// Scan implements the Scanner interface for database reading
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case string:
		parsed, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into TimeOnly", value)
	}
}

// Value implements the Valuer interface for database writing
func (t TimeOnly) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time.Format("15:04:05"), nil
}

// MarshalJSON implements json.Marshaler
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format("15:04:05"))), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}
	// Remove quotes
	str = str[1 : len(str)-1]
	parsed, err := time.Parse("15:04:05", str)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// BeforeCreate hook to generate UUID before creating
func (r *Restaurant) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Restaurant model
func (Restaurant) TableName() string {
	return "restaurants"
}
