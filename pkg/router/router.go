package router

import (
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/handlers"
	"github.com/HMZElidrissi/hotel-reservation-system-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func SetupRouter(db *mongo.Client) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")

	// Public routes
	{
		api.POST("/register", func(c *gin.Context) { handlers.Register(c, db) })
		api.POST("/login", func(c *gin.Context) { handlers.Login(c, db) })
	}

	// Secure routes
	admin := api.Group("/admin")
	admin.Use(middleware.Auth(), middleware.Role("admin"))
	{
		admin.POST("/rooms", func(c *gin.Context) { handlers.CreateRoom(c, db) })
		admin.GET("/rooms", func(c *gin.Context) { handlers.GetRooms(c, db) })
		admin.GET("/rooms/:id", func(c *gin.Context) { handlers.GetRoom(c, db) })
		admin.PUT("/rooms/:id", func(c *gin.Context) { handlers.UpdateRoom(c, db) })
		admin.DELETE("/rooms/:id", func(c *gin.Context) { handlers.DeleteRoom(c, db) })
		admin.POST("/reservations", func(c *gin.Context) { handlers.CreateReservation(c, db) })
		admin.GET("/reservations", func(c *gin.Context) { handlers.GetReservations(c, db) })
		admin.GET("/reservations/:id", func(c *gin.Context) { handlers.GetReservation(c, db) })
		admin.PUT("/reservations/:id", func(c *gin.Context) { handlers.UpdateReservation(c, db) })
		admin.DELETE("/reservations/:id", func(c *gin.Context) { handlers.DeleteReservation(c, db) })
	}

	customer := api.Group("/customer")
	customer.Use(middleware.Auth(), middleware.Role("customer"))
	{
		customer.POST("/rooms/:id", func(c *gin.Context) { handlers.BookRoom(c, db) })
		customer.GET("/reservations", func(c *gin.Context) { handlers.GetCustomerReservations(c, db) })
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return router
}
