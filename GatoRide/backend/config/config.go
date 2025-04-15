package config

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv" // Import dotenv package
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ✅ Connect to MongoDB
func ConnectDB() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: No .env file found, using system environment variables.")
	}

	// Get MongoDB URI & DB name
	mongoURI := "mongodb+srv://smohammad:iCPpHajYSS6WNsGa@cluster0.yeggn.mongodb.net/gatorides?retryWrites=true&w=majority&appName=Cluster0"
	dbName := "gatorides"

	if mongoURI == "" || dbName == "" {
		log.Fatal("❌ MONGO_URI or DB_NAME is missing. Ensure .env is set up correctly.")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB Connection Failed:", err)
	}

	// Set global DB variable
	DB = client.Database(dbName)
	fmt.Println("✅ Connected to MongoDB:", dbName)
}

// ✅ Fetch Collection
func GetCollection(name string) *mongo.Collection {
	if DB == nil {
		log.Println("❌ Error: Database is not connected.")
		return nil
	}
	return DB.Collection(name)
}
