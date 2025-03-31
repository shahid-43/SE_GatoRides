package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserProfile(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userCollection := config.GetCollection("users")
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	// Remove sensitive fields
	user.Password = ""
	user.VerificationToken = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetUserRides(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	rideCollection := config.GetCollection("rides")

	// Rides offered
	ridesOfferedCursor, err1 := rideCollection.Find(context.TODO(), bson.M{"driver_id": userID})
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch offered rides"})
		return
	}
	var ridesOffered []models.Ride
	_ = ridesOfferedCursor.All(context.TODO(), &ridesOffered)

	// Rides taken
	ridesTakenCursor, err2 := rideCollection.Find(context.TODO(), bson.M{"passenger_ids": userID})
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch taken rides"})
		return
	}
	var ridesTaken []models.Ride
	_ = ridesTakenCursor.All(context.TODO(), &ridesTaken)

	c.JSON(http.StatusOK, gin.H{
		"rides_offered": ridesOffered,
		"rides_taken":   ridesTaken,
	})
}
