package database

import (
	"context"
	"github.com/HMZElidrissi/hotel-reservation-system-api/configs"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func InitDB() *mongo.Client {
	// Load configuration
	configs.LoadConfig()

	// Get configuration values
	mongoURI := viper.GetString("MONGO_DB_URI")

	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	return client
}
