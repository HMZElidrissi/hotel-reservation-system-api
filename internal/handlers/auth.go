package handlers

import (
	"context"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Register(c *gin.Context, client *mongo.Client) {
	// Create a new user model
	var user = models.User{}

	// Bind the JSON body to the user model (meaning the user model now has the data from the request)
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while hashing the password"})
	}

	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()

	// Get the collection
	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("users")

	// Create a context with a timeout (i.e. the context will be cancelled after 5 seconds if the operation is not completed)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Defer the cancellation of the context (i.e. the context will be cancelled when the function ends)
	defer cancel()

	// TODO: Check if the user already exists

	// Insert the user into the collection
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while inserting the user"})
		return
	}
	// Return a success message
	c.JSON(http.StatusCreated, gin.H{"message": "your account has been created"})
}
