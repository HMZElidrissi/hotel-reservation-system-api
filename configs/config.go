package configs

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

func LoadConfig() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Set default values for environment variables
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("MONGO_DB_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGO_DB_NAME", "hotel-reservation")

	// Set environment variables
	viper.AutomaticEnv()

	// Set configuration file name and path
	/*
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Error reading config file: %v", err)
		}
	*/
}
