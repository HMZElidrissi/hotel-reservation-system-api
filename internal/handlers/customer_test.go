package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/models"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	viper.Set("MONGO_DB_NAME", "test_db")
	viper.Set("JWT_SECRET", "your_secret_key")
}

func createJWT(t *testing.T, id primitive.ObjectID, email, role string) string {
	token, err := utils.GenerateJWT(id, email, role)
	assert.NoError(t, err)
	return token
}

func authMiddleware(t *testing.T, id primitive.ObjectID, email, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := createJWT(t, id, email, role)
		c.Request.Header.Set("Authorization", "Bearer "+token)
		claims, err := utils.ParseToken(token)
		assert.NoError(t, err)
		c.Set("user", claims)
		c.Next()
	}
}

func TestBookRoom(t *testing.T) {
	client, cleanup := setupTestDB()
	defer cleanup()

	r := gin.Default()
	r.Use(authMiddleware(t, primitive.NewObjectID(), "customer@example.com", "customer"))
	r.POST("/customer/rooms/:id", func(c *gin.Context) {
		BookRoom(c, client)
	})

	reqBody := `{
		"check_in": "2023-06-10",
		"check_out": "2023-06-15",
		"status": "confirmed"
	}`
	roomID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("POST", "/customer/rooms/"+roomID, bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetCustomerReservations(t *testing.T) {
	client, cleanup := setupTestDB()
	defer cleanup()

	userID := primitive.NewObjectID()
	r := gin.Default()
	r.Use(authMiddleware(t, userID, "customer@example.com", "customer"))

	// Insert a reservation for the user
	collection := client.Database(viper.GetString("MONGO_DB_NAME")).Collection("reservations")
	reservation := models.Reservation{
		ID:       primitive.NewObjectID(),
		UserID:   userID,
		RoomID:   primitive.NewObjectID(),
		CheckIn:  "2023-06-10",
		CheckOut: "2023-06-15",
		Status:   "confirmed",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, reservation)
	assert.NoError(t, err)

	r.GET("/customer/reservations", func(c *gin.Context) {
		GetCustomerReservations(c, client)
	})

	req, _ := http.NewRequest("GET", "/customer/reservations", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var reservations []models.Reservation
	err = json.Unmarshal(w.Body.Bytes(), &reservations)
	assert.NoError(t, err)
	assert.NotEmpty(t, reservations)
}
