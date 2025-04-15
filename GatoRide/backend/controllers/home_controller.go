package controllers

import (
	"backend/config"
	"backend/models"
	"backend/tests/mocks"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HomeHandler(c *gin.Context) {
	// Extract userID from context
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

	// Get users collection
	userCollection := config.GetCollection("users")

	// Query MongoDB using ObjectID
	var user models.User
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}

	// Check if user location is set
	fmt.Println("User Data:", user)
	if user.Location.Latitude == 0 && user.Location.Longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User location not set"})
		return
	}

	// Fetch nearby rides using geospatial query
	rides, err := FetchNearbyRides(user.Location.Latitude, user.Location.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rides", "details": err.Error()})
		return
	}

	// Return rides
	c.JSON(http.StatusOK, gin.H{"rides": rides})
}

func FetchNearbyRides(lat float64, lon float64) ([]models.Ride, error) {
	collection := config.GetCollection("rides")
	radius := 1110.0 // 10 km search radius

	// Approximate degrees per kilometer (1° latitude ≈ 111 km)
	degreeOffset := radius / 111

	// Define bounding box filter
	filter := bson.M{
		"pickup.latitude":  bson.M{"$gte": lat - degreeOffset, "$lte": lat + degreeOffset},
		"pickup.longitude": bson.M{"$gte": lon - degreeOffset, "$lte": lon + degreeOffset},
		"status":           "open",
	}

	// Query MongoDB for nearby rides
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decode rides
	var rides []models.Ride
	for cursor.Next(context.TODO()) {
		var ride models.Ride
		if err := cursor.Decode(&ride); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, nil
}

func HomeHandlerWithMocks(userCollection mocks.CollectionInterfaces, ridesCollection mocks.CollectionInterfaces) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		var user models.User
		err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if user.Location.Latitude == 0 && user.Location.Longitude == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User location not set"})
			return
		}

		// ✅ Use the new mock-friendly function
		rides, err := FetchNearbyRidesWithMocks(ridesCollection, user.Location.Latitude, user.Location.Longitude)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rides"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"rides": rides})
	}
}

func FetchNearbyRidesWithMocks(ridesCollection mocks.CollectionInterfaces, lat float64, lon float64) ([]models.Ride, error) {
	radius := 10.0               // 10 km search radius
	degreeOffset := radius / 111 // Approximate conversion

	filter := bson.M{
		"pickup.latitude":  bson.M{"$gte": lat - degreeOffset, "$lte": lat + degreeOffset},
		"pickup.longitude": bson.M{"$gte": lon - degreeOffset, "$lte": lon + degreeOffset},
		"status":           "open",
	}

	// Query the mock collection
	cursor, err := ridesCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decode rides
	var rides []models.Ride
	for cursor.Next(context.TODO()) {
		var ride models.Ride
		if err := cursor.Decode(&ride); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, nil
}
