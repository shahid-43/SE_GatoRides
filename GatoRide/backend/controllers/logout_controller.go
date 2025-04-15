package controllers

import (
	"backend/config"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// LogOut handles user logout by invalidating the current session
func LogOut(c *gin.Context) {
	fmt.Println("LogOut func entered")
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization token required"})
		return
	}

	// Remove "Bearer " prefix if present
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the sessions collection
	sessionCollection := config.GetCollection("sessions")
	if sessionCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	// Delete the session from the database
	result, err := sessionCollection.DeleteOne(ctx, bson.M{"token": tokenString})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// Check if a session was actually deleted
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
