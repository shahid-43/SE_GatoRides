package controllers

import (
	"backend/config"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateUserProfile handles updating user profile information
func UpdateUserProfile(c *gin.Context) {
	// Extract userID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert userID to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse request body
	var updateData struct {
		Name     string          `json:"name"`
		Username string          `json:"username"`
		Location LocationRequest `json:"location"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get the users collection
	userCollection := config.GetCollection("users")
	if userCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	// Check if username already exists (if username is being updated)
	if updateData.Username != "" {
		// First check if the username is already taken by another user
		var existingUser map[string]interface{}
		err := userCollection.FindOne(
			context.TODO(),
			bson.M{
				"username": updateData.Username,
				"_id":      bson.M{"$ne": userID}, // Exclude current user
			},
		).Decode(&existingUser)

		// If a user was found, the username is taken
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
	}

	// Create update document
	updateFields := bson.M{}

	// Only add fields that are provided
	if updateData.Name != "" {
		updateFields["name"] = updateData.Name
	}
	if updateData.Username != "" {
		updateFields["username"] = updateData.Username
	}
	if (updateData.Location != LocationRequest{}) {
		updateFields["location"] = bson.M{
			"latitude":  updateData.Location.Latitude,
			"longitude": updateData.Location.Longitude,
			"address":   updateData.Location.Address,
		}
	}

	// If no fields were provided, return error
	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Update the user
	updateResult, err := userCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userID},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Check if any documents were updated
	if updateResult.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetch updated user to return in response
	var updatedUser map[string]interface{}
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&updatedUser)
	if err != nil {
		// If we can't fetch the updated user, still return success
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
		return
	}

	// Don't return sensitive fields
	delete(updatedUser, "password")
	delete(updatedUser, "verification_token")

	// Return success with updated user data
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

// LocationRequest is a helper struct for location updates
type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}
