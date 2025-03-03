package controllers_test

import (
	"backend/controllers"
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock ride data
var mockRide = models.Ride{
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
}

// Mock function for FetchRideFeedData
func mockFetchRideFeedData(lat float64, lon float64, date time.Time) ([]models.Ride, error) {
	fmt.Printf("Mock FetchRideFeedData called with lat: %f, lon: %f, date: %s\n", lat, lon, date)
	return []models.Ride{mockRide}, nil
}

// Test FetchRideFeed API
func TestFetchRideFeed(t *testing.T) {
	// Set Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new router and register the endpoint
	router := gin.Default()
	router.POST("/rides/feed", controllers.FetchRideFeed)

	// Mock request body
	requestBody, _ := json.Marshal(map[string]interface{}{
		"latitude":  40.7128,
		"longitude": -74.0060,
		"date":      time.Now().Format(time.RFC3339),
	})

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/rides/feed", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Response should be a valid JSON")

	rides, ok := response["rides"].([]interface{})
	assert.True(t, ok, "Rides should be an array")
	assert.Greater(t, len(rides), 0, "At least one ride should be returned")
}
