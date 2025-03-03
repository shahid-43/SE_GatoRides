package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// GetCollection is a function variable that can be overridden in tests
var GetCollection = func(name string) *mongo.Collection {
	if DB == nil {
		return nil
	}
	return DB.Collection(name)
}

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(os.Getenv("DB_NAME"))
	fmt.Println("Connected to MongoDB!")
}
