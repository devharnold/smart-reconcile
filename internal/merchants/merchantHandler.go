package merchants

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *UserService
}

func NewHandler(s *UserService) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	var body SignUpRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.Service.RegisterUser(
		c.Request.Context(),
		body.FirstName,
		body.LastName,
		body.UserEmail,
		body.PhoneNumber,
		body.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H {
		"id": user.UserID,
		"first_name": user.FirstName,
		"last_name": user.LastName,
		"email": user.Email,
		"phone_number": user.PhoneNumber,
		"token": user.Token,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var body LoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.Service.LoginUser(
		c.Request.Context(),
		body.UserEmail,
		body.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": user.Token,
	})
}

func (h *Handler) SearchUser(c *gin.Context) {
	var body SearchUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.Service.GetUserByEmail(
		c.Request.Context(),
		body.UserEmail,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.UserID,
		"first_name": user.FirstName,
		"last_name": user.LastName,
		"email": user.Email,
	})
}