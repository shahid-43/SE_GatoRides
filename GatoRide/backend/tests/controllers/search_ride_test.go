package controllers_test

import (
	"backend/config"
	"backend/controllers"
	"backend/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupSearchRideRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/user/search-ride", controllers.SearchRides)
	return router
}

func TestSearchRides_Success(t *testing.T) {
	router := setupSearchRideRouter()

	// Arrange: insert a ride into DB
	ride := models.Ride{
		ID:        primitive.NewObjectID(),
		DriverID:  primitive.NewObjectID(),
		Pickup:    models.Location{Latitude: 40.7128, Longitude: -74.0060, Address: "NYC"},
		Dropoff:   models.Location{Latitude: 40.7306, Longitude: -73.9352, Address: "Brooklyn"},
		Status:    models.StatusOpen,
		Price:     20,
		Seats:     3,
		Date:      time.Now(),
		CreatedAt: time.Now(),
	}
	collection := config.GetCollection("rides")
	_, err := collection.InsertOne(context.TODO(), ride)
	assert.NoError(t, err)

	// Prepare request payload
	requestBody, _ := json.Marshal(map[string]interface{}{
		"from": map[string]interface{}{
			"latitude":  ride.Pickup.Latitude,
			"longitude": ride.Pickup.Longitude,
		},
		"to": map[string]interface{}{
			"latitude":  ride.Dropoff.Latitude,
			"longitude": ride.Dropoff.Longitude,
		},
		"date":  time.Now().Format("2006-01-02"),
		"seats": 2,
	})

	req, _ := http.NewRequest("POST", "/user/search-ride", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string][]models.Ride
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response["rides"]), 1)

	// Clean up
	_, _ = collection.DeleteOne(context.TODO(), bson.M{"_id": ride.ID})
}

func TestSearchRides_InvalidDate(t *testing.T) {
	router := setupSearchRideRouter()

	requestBody := []byte(`{
		"from": {"latitude": 40.7128, "longitude": -74.0060},
		"to": {"latitude": 40.7306, "longitude": -73.9352},
		"date": "31-03-2025",
		"seats": 1
	}`)

	req, _ := http.NewRequest("POST", "/user/search-ride", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchRides_InvalidRequest(t *testing.T) {
	router := setupSearchRideRouter()

	req, _ := http.NewRequest("POST", "/user/search-ride", bytes.NewBuffer([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
