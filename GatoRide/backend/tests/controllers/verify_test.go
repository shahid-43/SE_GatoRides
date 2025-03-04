package controllers_test

import (
	"backend/config"
	"backend/controllers"
	"backend/models"
	"backend/utils"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestVerifyEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/verify-email", controllers.VerifyEmail)

	t.Run("Successful verification", func(t *testing.T) {
		email := "verify@example.com"
		verificationToken, _ := utils.GenerateVerificationToken(email)

		testUser := models.User{
			Username:          "verifytest",
			Email:             email,
			Password:          "Password123!",
			IsVerified:        false,
			VerificationToken: verificationToken,
		}

		cleanupTestUser(t, testUser.Email, testUser.Username)
		defer cleanupTestUser(t, testUser.Email, testUser.Username)

		collection := config.GetCollection("users")
		hashedPassword, _ := utils.HashPassword(testUser.Password)
		testUser.Password = hashedPassword
		testUser.ID = primitive.NewObjectID()
		_, _ = collection.InsertOne(context.TODO(), testUser)

		req, _ := http.NewRequest("GET", "/verify-email?token="+verificationToken, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
