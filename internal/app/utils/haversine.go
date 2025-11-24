package utils

import (
	"math"

	"github.com/google/uuid"
)

// CalculateDistance calculates the distance between two points using Haversine formula
// Returns distance in kilometers
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371.0 // Earth's radius in kilometers

	// Convert degrees to radians
	lat1Rad := degreesToRadians(lat1)
	lng1Rad := degreesToRadians(lng1)
	lat2Rad := degreesToRadians(lat2)
	lng2Rad := degreesToRadians(lng2)

	// Calculate differences
	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}

// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// RestaurantDistance represents a restaurant with its distance from a point
type RestaurantDistance struct {
	RestaurantID uuid.UUID
	Distance     float64
}
