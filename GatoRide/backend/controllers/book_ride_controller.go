// book_controller.go (updated with booking alerts)

package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookRide - Request a ride by adding a booking alert
func BookRide(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rideID := c.Query("ride_id")
	if rideID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ride ID"})
		return
	}

	rideObjectID, err := primitive.ObjectIDFromHex(rideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	rideCollection := config.GetCollection("rides")
	var ride models.Ride
	err = rideCollection.FindOne(context.TODO(), bson.M{"_id": rideObjectID}).Decode(&ride)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found"})
		return
	}

	if ride.Seats <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No seats available"})
		return
	}

	// Add booking alert to a new collection
	booking := models.BookingAlert{
		ID:        primitive.NewObjectID(),
		RideID:    ride.ID,
		Passenger: userIDStr.(string),
		DriverID:  ride.DriverID.Hex(),
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	bookingCollection := config.GetCollection("booking_alerts")
	_, err = bookingCollection.InsertOne(context.TODO(), booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send booking request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking request sent to driver. Waiting for confirmation."})
}

// AcceptBooking - Driver accepts a booking request
func AcceptBooking(c *gin.Context) {
	driverIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bookingID := c.Query("booking_id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing booking ID"})
		return
	}

	bookingObjID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	bookingCollection := config.GetCollection("booking_alerts")
	rideCollection := config.GetCollection("rides")

	// Update booking status
	update := bson.M{"$set": bson.M{"status": "confirmed"}}
	filter := bson.M{"_id": bookingObjID, "driver_id": driverIDStr}
	result, err := bookingCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm booking or not authorized"})
		return
	}

	// Decrease available seats in the ride
	booking := models.BookingAlert{}
	_ = bookingCollection.FindOne(context.TODO(), bson.M{"_id": bookingObjID}).Decode(&booking)
	rideID := booking.RideID
	_, _ = rideCollection.UpdateOne(context.TODO(), bson.M{"_id": rideID}, bson.M{"$inc": bson.M{"seats": -1}})

	c.JSON(http.StatusOK, gin.H{"message": "Booking confirmed."})
}
