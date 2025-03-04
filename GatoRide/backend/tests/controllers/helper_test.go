package controllers_test

import (
	"backend/config"
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// Helper function to clean up test users
func cleanupTestUser(t *testing.T, email string, username string) {
	collection := config.GetCollection("users")
	filter := bson.M{}

	if email != "" && username != "" {
		filter = bson.M{"$or": []bson.M{
			{"email": email},
			{"username": username},
		}}
	} else if email != "" {
		filter = bson.M{"email": email}
	} else if username != "" {
		filter = bson.M{"username": username}
	}

	if len(filter) > 0 {
		_, _ = collection.DeleteMany(context.TODO(), filter)
	}
}

// Helper function to clean up user sessions
func cleanupUserSessions(t *testing.T, userID string) {
	if userID == "" {
		return
	}

	collection := config.GetCollection("sessions")
	_, _ = collection.DeleteMany(context.TODO(), bson.M{"user_id": userID})
}
