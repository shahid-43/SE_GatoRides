package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(email, verificationToken string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	to := []string{email}
	subject := "Verify Your Email"
	body := fmt.Sprintf("Click the link to verify your email: http://localhost:8080/verify-email?token=%s", verificationToken)

	message := []byte("Subject: " + subject + "\r\n" + "\r\n" + body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("SMTP Error:", err) // Debugging
		return err
	}

	fmt.Println("Verification email sent successfully to:", email)
	return nil
}