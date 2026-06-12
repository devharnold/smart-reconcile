package router

import (
	"github.com/devharnold/smart-reconcile/internal/auth"
	"github.com/devharnold/smart-reconcile/internal/authmiddleware"
	"github.com/devharnold/smart-reconcile/internal/merchants"
	"github.com/gin-gonic/gin"
)

func SetUpRouter(handler *merchants.Handler, jwtSvc auth.JWTService) *gin.Engine {

	r := gin.Default()

	r.POST("/signup", handler.SignUp)
	r.POST("/login", handler.Login)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuth(jwtSvc))
	{
		protected.POST("/search-user", handler.SearchUser)
	}

	return r
}
