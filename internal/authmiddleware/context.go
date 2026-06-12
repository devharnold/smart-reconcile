package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) string {
	return c.MustGet("user_id").(string)
}

func GetEmail(c *gin.Context) string {
	return c.MustGet("email").(string)
}

func GetRole(c *gin.Context) string {
	return c.MustGet("role").(string)
}
