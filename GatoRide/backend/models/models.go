package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name"`
	Email             string             `bson:"email"`
	Username          string             `bson:"username"`
	Password          string             `bson:"password"`
	IsVerified        bool               `bson:"is_verified"`
	VerificationToken string             `bson:"verification_token"`
}
