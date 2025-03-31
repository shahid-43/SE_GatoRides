package controllers_test

import (
	"backend/config"
	"backend/controllers"
	"backend/models"
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

const mockID = "65f2d9e2c0b2a2e6b6b3a1d9"

func setupProfileRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("userID", mockID)
		c.Next()
	})

	router.GET("/user/profile", controllers.GetUserProfile)
	router.GET("/user/profile/rides", controllers.GetUserRides)

	return router
}

func insertTestUser(t *testing.T) {
	collection := config.GetCollection("users")
	userID, _ := primitive.ObjectIDFromHex(mockID)
	testUser := models.User{
		ID:         userID,
		Name:       "Test User",
		Email:      "test@example.com",
		Username:   "testuser",
		Password:   "hashed-password",
		IsVerified: true,
	}
	_, _ = collection.InsertOne(context.TODO(), testUser)
}

func insertingTestRides(t *testing.T) {
	collection := config.GetCollection("rides")
	userID, _ := primitive.ObjectIDFromHex(mockID)

	ride1 := models.Ride{
		ID:        primitive.NewObjectID(),
		DriverID:  userID,
		Pickup:    models.Location{Latitude: 1.0, Longitude: 2.0, Address: "Start 1"},
		Dropoff:   models.Location{Latitude: 3.0, Longitude: 4.0, Address: "End 1"},
		Price:     10,
		Status:    "open",
		CreatedAt: time.Now(),
	}

	ride2 := models.Ride{
		ID:           primitive.NewObjectID(),
		PassengerIDs: []primitive.ObjectID{userID},
		Pickup:       models.Location{Latitude: 5.0, Longitude: 6.0, Address: "Start 2"},
		Dropoff:      models.Location{Latitude: 7.0, Longitude: 8.0, Address: "End 2"},
		Price:        20,
		Status:       "completed",
		CreatedAt:    time.Now(),
	}

	_, _ = collection.InsertMany(context.TODO(), []interface{}{ride1, ride2})
}

func cleanupTestUserAndRides(t *testing.T) {
	userID, _ := primitive.ObjectIDFromHex(mockID)
	_ = config.GetCollection("users").FindOneAndDelete(context.TODO(), bson.M{"_id": userID})
	_, _ = config.GetCollection("rides").DeleteMany(context.TODO(), bson.M{
		"$or": []bson.M{
			{"driver_id": userID},
			{"passenger_ids": userID},
		},
	})
}

func TestGetUserProfile(t *testing.T) {
	insertTestUser(t)
	defer cleanupTestUserAndRides(t)

	router := setupProfileRouter()
	req, _ := http.NewRequest("GET", "/user/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "user")
}

func TestGetUserRides(t *testing.T) {
	insertTestUser(t)
	insertingTestRides(t)
	defer cleanupTestUserAndRides(t)

	router := setupProfileRouter()
	req, _ := http.NewRequest("GET", "/user/profile/rides", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "rides_offered")
	assert.Contains(t, response, "rides_taken")
}
