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

func CreateRoom(c *gin.Context, client *mongo.Client) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room.ID = primitive.NewObjectID()
	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("rooms")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating the room"})
		return
	}

	c.JSON(http.StatusCreated, room)
}

func GetRooms(c *gin.Context, client *mongo.Client) {
	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("rooms")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the rooms"})
		return
	}
	defer cursor.Close(ctx)

	var rooms []models.Room
	if err = cursor.All(ctx, &rooms); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while decoding the rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func GetRoom(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id) // Convert the string ID to an ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("rooms")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var room models.Room
	if err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&room); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func UpdateRoom(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("rooms")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"number": room.Number,
			"type":   room.Type,
			"price":  room.Price,
			"status": room.Status,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating the room"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func DeleteRoom(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("rooms")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleting the room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "room has been deleted"})
}
