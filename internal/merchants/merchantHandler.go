package merchants

import (
	"net/http"

	middleware "github.com/devharnold/smart-reconcile/internal/authmiddleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *MerchantService
}

func NewMerchantHandler(s *MerchantService) *Handler {
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

	merchant, err := h.Service.RegisterMerchant(
		c.Request.Context(),
		body.BusinessName,
		body.Email,
		body.PhoneNumber,
		body.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            merchant.UserID,
		"business_name": merchant.BusinessName,
		"email":         merchant.Email,
		"phone_number":  merchant.PhoneNumber,
		"token":         merchant.Token,
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

	merchant, err := h.Service.MerchantLogin(
		c.Request.Context(),
		body.Email,
		body.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": merchant.Token,
	})
}

func (h *Handler) SearchUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	if role != "merchant" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})
		return
	}

	var body SearchUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	merchant, err := h.Service.GetUserByEmail(
		c.Request.Context(),
		body.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requested_by": userID,
		"id":           merchant.UserID,
		"first_name":   merchant.BusinessName,
		"email":        merchant.Email,
	})
}
