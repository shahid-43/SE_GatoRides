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

// BookRide allows a user to book a ride by ride_id
func BookRide(c *gin.Context) {
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

	type Request struct {
		RideID string `json:"ride_id"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	rideID, err := primitive.ObjectIDFromHex(req.RideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	ridesCol := config.GetCollection("rides")

	// Check if ride exists
	var ride models.Ride
	err = ridesCol.FindOne(context.TODO(), bson.M{"_id": rideID}).Decode(&ride)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found"})
		return
	}

	// Check if already booked
	for _, passengerID := range ride.PassengerIDs {
		if passengerID == userID {
			c.JSON(http.StatusConflict, gin.H{"error": "Already booked this ride"})
			return
		}
	}

	// Update ride with new passenger
	update := bson.M{"$push": bson.M{"passenger_ids": userID}}
	_, err = ridesCol.UpdateOne(context.TODO(), bson.M{"_id": rideID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book ride"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ride booked successfully"})
}
