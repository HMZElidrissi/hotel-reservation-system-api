package handlers

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	viper.Set("MONGO_DB_NAME", "test_db")
	viper.Set("JWT_SECRET", "your_secret_key")
}

func setupTestDB() (*mongo.Client, func()) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Clean up function
	return client, func() {
		client.Database(viper.GetString("MONGO_DB_NAME")).Drop(context.Background())
		client.Disconnect(context.Background())
	}
}

func TestRegister(t *testing.T) {
	client, cleanup := setupTestDB()
	defer cleanup()

	r := gin.Default()
	r.POST("/register", func(c *gin.Context) {
		Register(c, client)
	})

	reqBody := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password": "password123"
	}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "your account has been created")
}

func TestRegisterUserAlreadyExists(t *testing.T) {
	client, cleanup := setupTestDB()
	defer cleanup()

	r := gin.Default()
	r.POST("/register", func(c *gin.Context) {
		Register(c, client)
	})

	reqBody := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password": "password123"
	}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "user already exists")
}

func TestLogin(t *testing.T) {
	client, cleanup := setupTestDB()
	defer cleanup()

	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		Login(c, client)
	})

	// Register a user first
	registerReqBody := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password": "password123"
	}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(registerReqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.POST("/register", func(c *gin.Context) {
		Register(c, client)
	})
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Attempt to login
	loginReqBody := `{
		"email": "john@example.com",
		"password": "password123"
	}`
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(loginReqBody)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLoginWrongPassword(t *testing.T) {
	client, cleanup := setupTestDB()
	defer cleanup()

	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		Login(c, client)
	})

	registerReqBody := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password": "password123"
	}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(registerReqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.POST("/register", func(c *gin.Context) {
		Register(c, client)
	})
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	loginReqBody := `{
		"email": "john@example.com",
		"password": "wrongpassword"
	}`
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(loginReqBody)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "wrong password")
}
