package handlers

import (
	"context"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func BookRoom(c *gin.Context, client *mongo.Client) {
	var reservation models.Reservation
	room, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	reservation.ID = primitive.NewObjectID()
	reservation.UserID = userID
	reservation.RoomID = room
	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("reservations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, reservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating the reservation"})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

func GetCustomerReservations(c *gin.Context, client *mongo.Client) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("reservations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the reservations"})
		return
	}
	defer cursor.Close(ctx)

	var reservations []models.Reservation
	if err = cursor.All(ctx, &reservations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while decoding the reservations"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}
