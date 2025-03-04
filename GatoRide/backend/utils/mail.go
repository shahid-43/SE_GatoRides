package utils

import (
	"fmt"
	"net/smtp"
)

var SendEmailFunc = SendVerificationEmail

func SendVerificationEmail(email string, verificationToken string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	from := "shahidshareef4457@gmail.com"
	password := "ldvqpjhklisopenq"

	to := []string{email}
	subject := "Verify Your Email"
	body := fmt.Sprintf("Click the link to verify your email: http://localhost:5001/verify-email?token=%s", verificationToken)

	message := []byte("Subject: " + subject + "\r\n" + "\r\n" + body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("SMTP Error:", err) // Debugging
		return err
	}

	fmt.Println("Verification email sent successfully to:", email)
	return nil
}
