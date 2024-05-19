package router

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func SetupRouter(db *mongo.Client) *gin.Engine {
	router := gin.Default()

	// Inject the database instance into the context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	/*// Public routes
	r.POST("/login", func(c *gin.Context) { handlers.Login(c, session) })

	// Secure routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		admin.POST("/users", func(c *gin.Context) { handlers.CreateUser(c, session) })
	}

	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("user"))
	{
		user.GET("/users/:id/posts", func(c *gin.Context) { handlers.GetUserWithPosts(c, session) })
		user.POST("/posts", func(c *gin.Context) { handlers.CreatePost(c, session) })
	}*/

	// Define the routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return router
}
