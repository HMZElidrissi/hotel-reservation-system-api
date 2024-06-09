package handlers

import (
	"context"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/models"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type RegisterRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `json:"first_name" binding:"required"`
	LastName  string             `json:"last_name" binding:"required"`
	Email     string             `json:"email" binding:"required,email"`
	Password  string             `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context, client *mongo.Client) {
	var req RegisterRequest

	// Bind the JSON body to the user model (meaning the user model now has the data from the request)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while hashing the password"})
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "customer",
	}

	// Get the collection
	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("users")

	// Create a context with a timeout (i.e. the context will be cancelled after 5 seconds if the operation is not completed)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Defer the cancellation of the context (i.e. the context will be cancelled when the function ends)
	defer cancel()

	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// Insert the user into the collection
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while inserting the user"})
		return
	}

	var storedUser models.User
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&storedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving the stored user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "your account has been created"})
}

func Login(c *gin.Context, client *mongo.Client) {
	var credentials LoginRequest

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&user) // Find the user by email and decode the result into the user model, decode it means that the result will be stored in the user model
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while generating the token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
