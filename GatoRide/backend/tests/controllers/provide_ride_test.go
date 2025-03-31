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

// 🔹 Mock user ID for authentication
const mockUserID = "65f2d9e2c0b2a2e6b6b3a1d9"

// ✅ **Setup test router with middleware**
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Apply authentication middleware
	router.Use(func(c *gin.Context) {
		c.Set("userID", mockUserID)
		c.Next()
	})

	router.POST("/user/provide-ride", controllers.ProvideRide)
	return router
}

// ✅ **Test: Successful Ride Creation**
func TestProvideRide_Success(t *testing.T) {
	router := setupRouter()

	requestBody, _ := json.Marshal(map[string]interface{}{
		"pickup": map[string]interface{}{
			"latitude":  40.7128,
			"longitude": -74.0060,
			"address":   "New York, NY",
		},
		"dropoff": map[string]interface{}{
			"latitude":  40.7306,
			"longitude": -73.9352,
			"address":   "Brooklyn, NY",
		},
		"price": 25.5,
		"seats": 3,
		"date":  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	})

	req, _ := http.NewRequest("POST", "/user/provide-ride", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "ride_id")

	rideIDStr := response["ride_id"].(string)
	rideID, err := primitive.ObjectIDFromHex(rideIDStr)
	assert.Nil(t, err)

	collection := config.GetCollection("rides")
	var createdRide models.Ride
	err = collection.FindOne(context.TODO(), bson.M{"_id": rideID}).Decode(&createdRide)
	assert.Nil(t, err)
	assert.Equal(t, createdRide.DriverID.Hex(), mockUserID)
	assert.Equal(t, 3, createdRide.Seats)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), createdRide.Date, time.Hour*24)

	cleanupTestRide(t, rideID)
}

// ✅ **Test: Unauthorized Request (No userID)**
func TestProvideRide_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/user/provide-ride", controllers.ProvideRide)

	req, _ := http.NewRequest("POST", "/user/provide-ride", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Unauthorized")
}

// ✅ **Test: Invalid Ride Data**
func TestProvideRide_InvalidData(t *testing.T) {
	router := setupRouter()

	requestBody, _ := json.Marshal(map[string]interface{}{
		"pickup":  map[string]interface{}{},
		"dropoff": map[string]interface{}{},
		"price":   0,
		"seats":   0,
		"date":    "",
	})

	req, _ := http.NewRequest("POST", "/user/provide-ride", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid ride data")
}

// 🛠 **Helper Function: Clean Up Test Ride**
func cleanupTestRide(t *testing.T, rideID primitive.ObjectID) {
	collection := config.GetCollection("rides")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": rideID})
	if err != nil {
		t.Logf("⚠️ Failed to clean up test ride: %v", err)
	}
}
