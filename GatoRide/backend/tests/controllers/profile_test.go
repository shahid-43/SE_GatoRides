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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUserProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Setup test user
	testUser := models.User{
		ID:       primitive.NewObjectID(),
		Name:     "Test User",
		Email:    "profile@example.com",
		Username: "profile_user",
		Password: "hashedpass",
	}
	collection := config.GetCollection("users")
	_, _ = collection.InsertOne(context.TODO(), testUser)
	defer cleanupTestUser(t, testUser.Email, testUser.Username)

	// Inject mock userID
	router.Use(func(c *gin.Context) {
		c.Set("userID", testUser.ID.Hex())
		c.Next()
	})
	router.GET("/user/profile", controllers.GetUserProfile)

	req, _ := http.NewRequest("GET", "/user/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "user")
}

func TestGetUserRides(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testUserID := primitive.NewObjectID()
	testUser := models.User{
		ID:       testUserID,
		Name:     "Ride Tester",
		Email:    "ridetest@example.com",
		Username: "ridetester",
		Password: "secret",
	}
	userCol := config.GetCollection("users")
	_, _ = userCol.InsertOne(context.TODO(), testUser)
	defer cleanupTestUser(t, testUser.Email, testUser.Username)

	ride := models.Ride{
		DriverID: testUserID,
		Status:   "open",
	}
	rideCol := config.GetCollection("rides")
	_, _ = rideCol.InsertOne(context.TODO(), ride)

	router.Use(func(c *gin.Context) {
		c.Set("userID", testUserID.Hex())
		c.Next()
	})
	router.GET("/user/profile/rides", controllers.GetUserRides)

	req, _ := http.NewRequest("GET", "/user/profile/rides", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "rides_offered")
	assert.Contains(t, response, "rides_taken")
}
