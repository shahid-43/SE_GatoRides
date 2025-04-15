package controllers

import (
	"backend/auth"
	"backend/config"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ✅ **Signup with Email Verification**
func Signup(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	collection := config.GetCollection("users")
	if collection == nil {
		log.Println("❌ Error: Users collection is nil! Ensure DB is connected.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
		return
	}
	if user.Email == "" || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, Username, and Password are required"})
		fmt.Println(user.Username)
		fmt.Println(user.Password)
		fmt.Println(user.Email)
		fmt.Println(user)
		return
	}
	// Check if email or username already exists
	var existingUser models.User
	errEmail := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	errUsername := collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)

	// Return precise error messages
	if errEmail == nil && errUsername == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email and Username already exist"})
		return
	} else if errEmail == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	} else if errUsername == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()
	user.IsVerified = false

	// Generate verification token
	verificationToken, _ := utils.GenerateVerificationToken(user.Email)
	user.VerificationToken = verificationToken

	// Insert user into DB
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Send verification email
	err = utils.SendVerificationEmail(user.Email, verificationToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		fmt.Println(user.Email)
		fmt.Println(verificationToken)
		fmt.Println(user)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created. Verify your email."})
}

// ✅ **Verify Email**
func VerifyEmail(c *gin.Context) {
	log.Println("Received email verification request.")

	// Extract token from query params
	token := c.Query("token")
	if token == "" {
		log.Println("Error: Missing verification token in request.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing verification token"})
		return
	}
	log.Println("Verification token received:", token)

	// Verify token and extract email
	email, err := utils.VerifyToken(token)
	if err != nil {
		log.Println("Error verifying token:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}
	log.Println("Token successfully verified. Email:", email)

	// Get the users collection
	collection := config.GetCollection("users")

	// Find user by email
	var user bson.M
	err = collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Println("Error: No user found with the given email:", email, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	log.Println("User found in database:", user)

	// Check if the user is already verified
	if user["is_verified"] == true {
		log.Println("User already verified:", email)
		c.JSON(http.StatusOK, gin.H{"message": "Email already verified."})
		return
	}

	// Update user verification status
	updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"email": email}, bson.M{"$set": bson.M{"is_verified": true}})
	if err != nil {
		log.Println("Error updating user verification status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying email"})
		return
	}

	// Log update result
	log.Println("Update Result:", updateResult)

	// Send successful response
	log.Println("Email successfully verified for:", email)
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully."})
}

// ✅ **Login User**
func Login(c *gin.Context) {
	// Decode request body
	var creds models.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	collection := config.GetCollection("users")
	var user models.User

	// Check if user exists in database
	err := collection.FindOne(context.TODO(), bson.M{"email": creds.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Ensure user is verified before allowing login
	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email before logging in."})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token expiration (matching JWT expiration: 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Create session document
	session := models.Session{
		UserID:    user.ID.Hex(),
		Token:     token,
		ExpiresAt: expirationTime,
	}

	// Store session in MongoDB
	sessionCollection := config.GetCollection("sessions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = sessionCollection.InsertOne(ctx, session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}

	// Send response with token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
