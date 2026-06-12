package merchants

type Merchants struct {
	UserID       int64  `json:"user_id"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Password     string `json:"-"`
	Token        string `json:"token,omitempty"`
}

type SignUpRequest struct {
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Password     string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SearchUserRequest struct {
	Email string `json:"email"`
}
