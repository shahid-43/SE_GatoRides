package controllers_test

import (
	"backend/config"
	"backend/controllers"
	"backend/models"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestLogout tests the logout functionality
func TestLogout(t *testing.T) {
	// Set up
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users/logout", controllers.LogOut)

	// Test successful logout
	t.Run("Successful logout", func(t *testing.T) {
		// Create a test session in the database
		sessionCollection := config.GetCollection("sessions")
		testToken := "test-token-12345"
		testSession := models.Session{
			UserID:    primitive.NewObjectID().Hex(),
			Token:     testToken,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		}

		_, err := sessionCollection.InsertOne(context.TODO(), testSession)
		assert.NoError(t, err)

		// Create request with the token
		req, _ := http.NewRequest("POST", "/users/logout", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify session was deleted
		var count int64
		count, err = sessionCollection.CountDocuments(context.TODO(), map[string]interface{}{
			"token": testToken,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	// Test logout with invalid token
	t.Run("Logout with invalid token", func(t *testing.T) {
		// Create request with an invalid token
		req, _ := http.NewRequest("POST", "/users/logout", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return not found for invalid token
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Test logout without token
	t.Run("Logout without token", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users/logout", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return bad request when no token provided
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
