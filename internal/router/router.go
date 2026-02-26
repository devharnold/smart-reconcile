package router

import (
	"github.com/gin-gonic/gin"
	"github.com/devharnold/smart-reconcile/internal/auth"
	"github.com/devharnold/smart-reconcile/internal/authmiddleware"
	"github.com/devharnold/smart-reconcile/internal/users"
)

// protected User Routes
func SetUpRouter(usersHandler *users.UserHandler, jwtSvc auth.JWTService) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public User Routes
	usersGroup := r.Group("/users")
	{
		usersGroup.POST("/signup", usersHandler.SignUp)
		usersGroup.POST("/login", usersHandler.Login)
	}

	// protected user routes
	protectedUsers := usersGroup.Group("/")
	protectedUsers.Use(
		middleware.JWTAuth(jwtSvc),
		middleware.RequireRole(auth.RoleAdmin, auth.RoleSupport),
	)
	{
		protectedUsers.POST("/search", usersHandler.SearchUser)
	}

	return r
}

// protected transaction routes
// func setUpTransactionRouter()