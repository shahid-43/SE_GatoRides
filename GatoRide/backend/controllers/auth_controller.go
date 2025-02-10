package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"backend/config"
	"backend/models"
	"backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection := config.GetCollection("users")

	// Check if email or username already exists
	var existingUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"email": user.Email}, {"username": user.Username}}}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "Email or Username already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()

	// Insert user into DB
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection := config.GetCollection("users")
	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"email": creds.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, _ := utils.GenerateJWT(user.Username)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
