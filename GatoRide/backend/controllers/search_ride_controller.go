package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const maxDistance = 0.05 // roughly ~5km range for latitude/longitude

type SearchRideRequest struct {
	From models.Location `json:"from" binding:"required"`
	To   models.Location `json:"to" binding:"required"`
	Date string          `json:"date" binding:"required"` // "YYYY-MM-DD"
}

func isNearby(loc1, loc2 models.Location) bool {
	return math.Abs(loc1.Latitude-loc2.Latitude) <= maxDistance &&
		math.Abs(loc1.Longitude-loc2.Longitude) <= maxDistance
}

func SearchRides(c *gin.Context) {
	var req SearchRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Parse date and get start and end of that day
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	startOfDay := date
	endOfDay := date.Add(24 * time.Hour)

	collection := config.GetCollection("rides")

	// Fetch all open rides within that day
	filter := bson.M{
		"status": "open",
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer cursor.Close(context.TODO())

	var matchingRides []models.Ride
	for cursor.Next(context.TODO()) {
		var ride models.Ride
		if err := cursor.Decode(&ride); err == nil {
			if isNearby(ride.Pickup, req.From) && isNearby(ride.Dropoff, req.To) {
				matchingRides = append(matchingRides, ride)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"rides": matchingRides})
}

// func GetUserProfile(c *gin.Context) {
// 	userIDStr, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	userCollection := config.GetCollection("users")
// 	var user models.User
// 	err := userCollection.FindOne(context.TODO(), bson.M{"_id": userIDStr}).Decode(&user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
// 		return
// 	}

// 	// Remove sensitive fields
// 	user.Password = ""
// 	user.VerificationToken = ""

// 	c.JSON(http.StatusOK, gin.H{"user": user})
// }

// func GetUserRides(c *gin.Context) {
// 	userIDStr, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	rideCollection := config.GetCollection("rides")

// 	// Rides offered
// 	ridesOfferedCursor, err1 := rideCollection.Find(context.TODO(), bson.M{"driver_id": userIDStr})
// 	if err1 != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch offered rides"})
// 		return
// 	}
// 	var ridesOffered []models.Ride
// 	_ = ridesOfferedCursor.All(context.TODO(), &ridesOffered)

// 	// Rides taken (assuming passenger_id field exists)
// 	ridesTakenCursor, err2 := rideCollection.Find(context.TODO(), bson.M{"passenger_id": userIDStr})
// 	if err2 != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch taken rides"})
// 		return
// 	}
// 	var ridesTaken []models.Ride
// 	_ = ridesTakenCursor.All(context.TODO(), &ridesTaken)

// 	c.JSON(http.StatusOK, gin.H{
// 		"rides_offered": ridesOffered,
// 		"rides_taken":   ridesTaken,
// 	})
// }
