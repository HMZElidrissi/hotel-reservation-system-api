package main

import (
	"context"
	"github.com/HMZElidrissi/hotel-reservation-system-api/configs"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/database"
	"github.com/HMZElidrissi/hotel-reservation-system-api/pkg/router"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var client *mongo.Client

func main() {
	configs.LoadConfig()

	client = database.InitDB()

	ctx := context.TODO()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect MongoDB client: %v", err)
		}
	}()

	serverPort := viper.GetString("SERVER_PORT")
	r := router.SetupRouter(client)
	err := r.Run(":" + serverPort)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
