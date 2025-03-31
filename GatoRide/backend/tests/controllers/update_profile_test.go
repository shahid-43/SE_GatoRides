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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateUserProfile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Add auth middleware mock and route
	router.POST("/user/update-profile", func(c *gin.Context) {
		// Mock the auth middleware by setting userID
		c.Set("userID", mockUserID)
		controllers.UpdateUserProfile(c)
	})

	// Test updating name and username
	t.Run("Update name and username", func(t *testing.T) {
		// Create a test user
		userID, err := primitive.ObjectIDFromHex(mockUserID)
		assert.NoError(t, err)

		// Insert test user
		testUser := models.User{
			ID:         userID,
			Name:       "Original Name",
			Email:      "test@example.com",
			Username:   "originalusername",
			Password:   "hashedpassword",
			IsVerified: true,
			Location: models.Location{
				Latitude:  0,
				Longitude: 0,
				Address:   "",
			},
		}

		// Clean up before and after test
		cleanupTestUser(t, testUser.Email, testUser.Username)
		defer cleanupTestUser(t, testUser.Email, testUser.Username)

		// Insert user into DB
		collection := config.GetCollection("users")
		_, err = collection.InsertOne(context.TODO(), testUser)
		assert.NoError(t, err)

		// Create update request
		updateData := map[string]interface{}{
			"name":     "Updated Name",
			"username": "updatedusername",
		}

		jsonData, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("POST", "/user/update-profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response contains success message
		assert.Contains(t, response, "message")
		assert.Contains(t, response["message"], "updated successfully")

		// Verify user contains updated data
		user, exists := response["user"].(map[string]interface{})
		assert.True(t, exists)
		assert.Equal(t, "Updated Name", user["name"])
		assert.Equal(t, "updatedusername", user["username"])

		// Verify database was updated
		var updatedUser models.User
		err = collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&updatedUser)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", updatedUser.Name)
		assert.Equal(t, "updatedusername", updatedUser.Username)
	})

	// Test updating location
	t.Run("Update location", func(t *testing.T) {
		// Create a test user
		userID, err := primitive.ObjectIDFromHex(mockUserID)
		assert.NoError(t, err)

		// Insert test user
		testUser := models.User{
			ID:         userID,
			Name:       "Location Test",
			Email:      "location@example.com",
			Username:   "locationtest",
			Password:   "hashedpassword",
			IsVerified: true,
			Location: models.Location{
				Latitude:  0,
				Longitude: 0,
				Address:   "",
			},
		}

		// Clean up before and after test
		cleanupTestUser(t, testUser.Email, testUser.Username)
		defer cleanupTestUser(t, testUser.Email, testUser.Username)

		// Insert user into DB
		collection := config.GetCollection("users")
		_, err = collection.InsertOne(context.TODO(), testUser)
		assert.NoError(t, err)

		// Create update request with location
		updateData := map[string]interface{}{
			"location": map[string]interface{}{
				"latitude":  40.7128,
				"longitude": -74.0060,
				"address":   "New York, NY",
			},
		}

		jsonData, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("POST", "/user/update-profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify database was updated
		var updatedUser models.User
		err = collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&updatedUser)
		assert.NoError(t, err)
		assert.Equal(t, 40.7128, updatedUser.Location.Latitude)
		assert.Equal(t, -74.0060, updatedUser.Location.Longitude)
		assert.Equal(t, "New York, NY", updatedUser.Location.Address)
	})

	// Test duplicate username
	t.Run("Duplicate username", func(t *testing.T) {
		// Create two test users
		userID, err := primitive.ObjectIDFromHex(mockUserID)
		assert.NoError(t, err)

		// User 1 (the user we'll update)
		testUser1 := models.User{
			ID:         userID,
			Name:       "Test User 1",
			Email:      "user1@example.com",
			Username:   "testuser1",
			Password:   "hashedpassword",
			IsVerified: true,
		}

		// User 2 (existing user with the username we want to take)
		testUser2 := models.User{
			ID:         primitive.NewObjectID(),
			Name:       "Test User 2",
			Email:      "user2@example.com",
			Username:   "existingusername", // we'll try to update User 1 to this
			Password:   "hashedpassword",
			IsVerified: true,
		}

		// Clean up before and after test
		cleanupTestUser(t, testUser1.Email, testUser1.Username)
		defer cleanupTestUser(t, testUser1.Email, testUser1.Username)
		cleanupTestUser(t, testUser2.Email, testUser2.Username)
		defer cleanupTestUser(t, testUser2.Email, testUser2.Username)

		// Insert users into DB
		collection := config.GetCollection("users")
		_, err = collection.InsertOne(context.TODO(), testUser1)
		assert.NoError(t, err)
		_, err = collection.InsertOne(context.TODO(), testUser2)
		assert.NoError(t, err)

		// Create update request with duplicate username
		updateData := map[string]interface{}{
			"username": "existingusername", // This username already exists for User 2
		}

		jsonData, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("POST", "/user/update-profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Username already exists")
	})

	// Test empty update
	t.Run("Empty update", func(t *testing.T) {
		// Create update request with empty data
		updateData := map[string]interface{}{}

		jsonData, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("POST", "/user/update-profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "No fields to update")
	})

	// Test invalid request format
	t.Run("Invalid request format", func(t *testing.T) {
		// Create invalid JSON request
		req, _ := http.NewRequest("POST", "/user/update-profile", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid request data")
	})
}
