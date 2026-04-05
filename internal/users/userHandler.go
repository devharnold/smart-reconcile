//todo: Migrate Users Handler to Lambda
// One architecture across the entire platform

package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *UserService
}

func NewUserHandler(s *UserService) *UserHandler {
	return &UserHandler{Service: s}
}

type SignUpRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	UserEmail   string `json:"user_email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type LoginRequest struct {
	UserEmail string `json:"user_email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type SearchUserRequest struct {
	UserEmail string `json:"user_email" binding:"required,email"`
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.RegisterUser(
		c.Request.Context(),
		req.FirstName,
		req.LastName,
		req.UserEmail,
		req.PhoneNumber,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           user.UserID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"token":        user.Token,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.LoginUser(
		c.Request.Context(),
		req.UserEmail,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": user.Token})
}

func (h *UserHandler) SearchUser(c *gin.Context) {
	var req SearchUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.GetUserByEmail(c.Request.Context(), req.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"id": user.UserID,
		"first_name": user.FirstName,
		"last_name": user.LastName,
		"email": user.Email,
	})
}