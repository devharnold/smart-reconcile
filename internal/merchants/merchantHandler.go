package merchants

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/devharnold/smart-reconcile/internal/authmiddleware"
	"github.com/devharnold/smart-reconcile/internal/http/response"

	"github.com/google/uuid"
)

type Handler struct {
	Service *MerchantService
}

func NewMerchantHandler(service *MerchantService) *Handler {
	return &Handler{
		Service: service,
	}
}

type SignUpResponse struct {
	ID           uuid.UUID `json:"id"`
	BusinessName string    `json:"business_name"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	Token        string    `json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SearchUserResponse struct {
	RequestedBy  string    `json:"requested_by"`
	ID           uuid.UUID `json:"id"`
	BusinessName string    `json:"business_name"`
	Email        string    `json:"email"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	merchant, err := h.Service.RegisterMerchant(
		r.Context(),
		req.BusinessName,
		req.Email,
		req.PhoneNumber,
		req.Password,
	)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusCreated, SignUpResponse{
		ID:           merchant.UserID,
		BusinessName: merchant.BusinessName,
		Email:        merchant.Email,
		PhoneNumber:  merchant.PhoneNumber,
		Token:        merchant.Token,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	merchant, err := h.Service.MerchantLogin(
		r.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, LoginResponse{
		Token: merchant.Token,
	})
}

func (h *Handler) SearchUser(w http.ResponseWriter, r *http.Request) {
	userID := authmiddleware.GetUserID(r.Context())
	role := authmiddleware.GetRole(r.Context())

	if role != "merchant" {
		response.JSON(w, http.StatusForbidden, ErrorResponse{
			Error: "access denied",
		})
		return
	}

	var req SearchUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	user, err := h.Service.GetUserByEmail(
		r.Context(),
		req.Email,
	)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, SearchUserResponse{
		RequestedBy:  userID,
		ID:           user.UserID,
		BusinessName: user.BusinessName,
		Email:        user.Email,
	})
}
