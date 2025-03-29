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

func TestHomeHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Add mock authentication middleware
	router.Use(func(c *gin.Context) {
		// Use the existing mock user ID
		c.Set("userID", mockUserID)
		c.Next()
	})

	router.GET("/home", controllers.HomeHandler)

	t.Run("Home handler with valid user and location", func(t *testing.T) {
		// Create user with location first
		userID, err := primitive.ObjectIDFromHex(mockUserID)
		assert.NoError(t, err)

		// Insert a test user with location
		insertTestUserWithLocation(t, userID)

		// Insert test rides
		insertTestRides(t)

		// Make request
		req, _ := http.NewRequest("GET", "/home", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response has rides
		rides, exists := response["rides"]
		assert.True(t, exists)
		assert.NotNil(t, rides)

		// Clean up
		cleanupTestUser(t, "home-test@example.com", "hometest")
		cleanupTestRides(t)
	})

	t.Run("Home handler with user without location", func(t *testing.T) {
		// Create user without location
		userID, err := primitive.ObjectIDFromHex(mockUserID)
		assert.NoError(t, err)

		// Insert a test user without location
		insertTestUserWithoutLocation(t, userID)

		// Make request
		req, _ := http.NewRequest("GET", "/home", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check status code - should be bad request when location not set
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Clean up
		cleanupTestUser(t, "home-test@example.com", "hometest")
	})
}

// Helper function to insert a test user with location
func insertTestUserWithLocation(t *testing.T, userID primitive.ObjectID) {
	collection := config.GetCollection("users")
	user := models.User{
		ID:       userID,
		Name:     "Home Test User",
		Email:    "home-test@example.com",
		Username: "hometest",
		Password: "password",
		Location: models.Location{
			Latitude:  40.7128,
			Longitude: -74.0060,
			Address:   "New York, NY",
		},
		IsVerified: true,
	}
	_, err := collection.InsertOne(context.TODO(), user)
	assert.NoError(t, err)
}

// Helper function to insert a test user without location
func insertTestUserWithoutLocation(t *testing.T, userID primitive.ObjectID) {
	collection := config.GetCollection("users")
	user := models.User{
		ID:         userID,
		Name:       "Home Test User",
		Email:      "home-test@example.com",
		Username:   "hometest",
		Password:   "password",
		IsVerified: true,
	}
	_, err := collection.InsertOne(context.TODO(), user)
	assert.NoError(t, err)
}

// Helper function to insert test rides
func insertTestRides(t *testing.T) {
	collection := config.GetCollection("rides")

	// Create rides near our test user location
	rides := []interface{}{
		models.Ride{
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
			Status:    models.StatusOpen,
			Price:     25.99,
			CreatedAt: time.Now(),
		},
		models.Ride{
			ID:       primitive.NewObjectID(),
			DriverID: primitive.NewObjectID(),
			Pickup: models.Location{
				Latitude:  40.7135,
				Longitude: -74.0046,
				Address:   "Manhattan, NY",
			},
			Dropoff: models.Location{
				Latitude:  40.6782,
				Longitude: -73.9442,
				Address:   "Brooklyn, NY",
			},
			Status:    models.StatusOpen,
			Price:     19.99,
			CreatedAt: time.Now(),
		},
	}

	_, err := collection.InsertMany(context.TODO(), rides)
	assert.NoError(t, err)
}

// Helper function to clean up test rides
func cleanupTestRides(t *testing.T) {
	collection := config.GetCollection("rides")
	_, err := collection.DeleteMany(context.TODO(), bson.M{"status": "open"})
	if err != nil {
		t.Logf("Error cleaning up test rides: %v", err)
	}
}
