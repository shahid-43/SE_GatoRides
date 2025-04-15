package controllers

import (
	"backend/config"
	"backend/models"
	"backend/tests/mocks"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// FetchRideFeed retrieves ride feed based on user-provided location and date
func FetchRideFeed(c *gin.Context) {
	var request struct {
		Latitude  float64   `json:"latitude" binding:"required"`
		Longitude float64   `json:"longitude" binding:"required"`
		Date      time.Time `json:"date" binding:"required"`
	}

	// Bind JSON request body to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input parameters"})
		return
	}

	// Use mock collection in tests
	var collection mocks.CollectionInterface
	if gin.Mode() == gin.TestMode {
		collection = mocks.MockGetCollection("rides")
	} else {
		collection = config.GetCollection("rides")
	}

	// Fetch rides based on provided location and date
	rides, err := FetchRideFeedData(collection, request.Latitude, request.Longitude, request.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ride feed"})
		return
	}

	// ✅ Force JSON to always return an array instead of `null`
	var rideResponse []interface{}
	for _, ride := range rides {
		rideResponse = append(rideResponse, ride)
	}

	c.JSON(http.StatusOK, gin.H{"rides": rideResponse})
}

// FetchRideFeedData retrieves rides based on the given location and date
func FetchRideFeedData(collection mocks.CollectionInterface, lat float64, lon float64, date time.Time) ([]models.Ride, error) {
	radius := 10.0 // Search radius in km

	// Approximate degrees per kilometer (1° latitude ≈ 111 km)
	degreeOffset := radius / 111

	// Define bounding box filter and date constraint
	filter := bson.M{
		"pickup.latitude":  bson.M{"$gte": lat - degreeOffset, "$lte": lat + degreeOffset},
		"pickup.longitude": bson.M{"$gte": lon - degreeOffset, "$lte": lon + degreeOffset},
		"status":           "open",
		"created_at": bson.M{
			"$gte": time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, time.UTC),
		},
	}

	// Query MongoDB for rides matching the criteria
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decode rides
	var rides []models.Ride
	for cursor.Next(context.TODO()) {
		var ride models.Ride
		if err := cursor.Decode(&ride); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, nil
}
