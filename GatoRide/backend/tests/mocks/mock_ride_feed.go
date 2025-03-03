package mocks

import (
	"backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockRides - Example ride data for tests
var MockRides = []models.Ride{
	{
		ID:       primitive.NewObjectID(),
		DriverID: primitive.NewObjectID(),
		Pickup: models.Location{
			Latitude:  40.7128,
			Longitude: -74.0060,
			Address:   "New York, NY",
		},
		Dropoff: models.Location{
			Latitude:  40.7306,
			Longitude: -73.9352,
			Address:   "Brooklyn, NY",
		},
		Status:    "open",
		Price:     25.5,
		CreatedAt: time.Now(),
	},
	{
		ID:       primitive.NewObjectID(),
		DriverID: primitive.NewObjectID(),
		Pickup: models.Location{
			Latitude:  37.7749,
			Longitude: -122.4194,
			Address:   "San Francisco, CA",
		},
		Dropoff: models.Location{
			Latitude:  37.8044,
			Longitude: -122.2711,
			Address:   "Oakland, CA",
		},
		Status:    "open",
		Price:     30.0,
		CreatedAt: time.Now(),
	},
}
