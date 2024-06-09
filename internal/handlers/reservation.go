package handlers

import (
	"context"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func CreateReservation(c *gin.Context, client *mongo.Client) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation.ID = primitive.NewObjectID()
	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("rooms")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, reservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating the reservation"})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

func GetReservations(c *gin.Context, client *mongo.Client) {
	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("reservations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
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

func GetReservation(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id) // Convert the string ID to an ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("reservations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reservation models.Reservation
	if err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&reservation); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func UpdateReservation(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("reservations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": reservation})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating the reservation"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func DeleteReservation(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("reservations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleting the reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reservation has been deleted"})
}
