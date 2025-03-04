package controllers_test

import (
	"backend/controllers"
	"backend/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// **Cleanup function to delete test users**

func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// ✅ Override SendEmailFunc to prevent actual email sending
	originalSendEmail := utils.SendEmailFunc
	utils.SendEmailFunc = func(email string, token string) error {
		fmt.Println("📧 Mock Email Sent to:", email)
		return nil
	}
	defer func() { utils.SendEmailFunc = originalSendEmail }() // Restore after test

	router.POST("/signup", controllers.Signup)

	t.Run("Successful signup", func(t *testing.T) {
		testEmail := "radime7497@lassora.com"
		testUsername := "testuser123"

		// **Clean up before the test**
		cleanupTestUser(t, testEmail, testUsername)

		requestBody, _ := json.Marshal(map[string]string{
			"email":    testEmail,
			"username": testUsername,
			"password": "Password123!",
			"name":     "Test User",
		})

		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		fmt.Println("Response:", w.Body.String()) // Debugging response

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response["message"], "User created")

		// **Clean up after the test**
		cleanupTestUser(t, testEmail, testUsername)
	})
}
