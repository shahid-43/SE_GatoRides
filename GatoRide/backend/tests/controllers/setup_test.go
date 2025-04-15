package controllers_test

import (
	"backend/config"
	"backend/utils"
	"fmt"
	"os"
	"testing"
)

// ✅ Ensure MongoDB is connected before running tests
func TestMain(m *testing.M) {
	// Connect to the database
	fmt.Println("🔄 Connecting to MongoDB for tests...")
	config.ConnectDB()

	// ✅ Override SendEmailFunc to prevent real emails
	originalSendEmail := utils.SendEmailFunc
	utils.SendEmailFunc = func(email string, token string) error {
		fmt.Println("📧 Mock Email Sent to:", email)
		return nil
	}
	defer func() { utils.SendEmailFunc = originalSendEmail }() // Restore after test

	// Run all tests
	exitCode := m.Run()

	// Exit with the proper exit code
	os.Exit(exitCode)
}
