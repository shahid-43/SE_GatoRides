package middlewares

import (
	"backend/auth"
	"backend/config"
	"backend/models"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthMiddleware validates JWT and checks session storage
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Extract the token from "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Verify session in MongoDB
		sessionCollection := config.GetCollection("sessions")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var session models.Session

		err = sessionCollection.FindOne(ctx, bson.M{"token": tokenString}).Decode(&session)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found. Please log in again."})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate session"})
			}
			c.Abort()
			return
		}

		// Check if the session is expired
		if time.Now().Unix() > session.ExpiresAt {
			fmt.Println("🔴 Session Expired!")
			fmt.Println("🔵 Current Time:", time.Now().Unix(), " | 🔴 Expiration Time:", session.ExpiresAt)
			sessionCollection.DeleteOne(ctx, bson.M{"token": tokenString})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired. Please log in again."})
			c.Abort()
			return
		}

		// Attach user info to context for later use
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
