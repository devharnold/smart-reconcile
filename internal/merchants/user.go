package users

type Users struct {
	UserID      int64  `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"-"`
	Token       string `json:"token,omitempty"`
}

type SignUpRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	UserEmail  string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password   string `json:"password"`
}

type LoginRequest struct {
	UserEmail string `json:"email"`
	Password  string `json:"password"`
}

type SearchUserRequest struct {
	UserEmail string `json:"email"`
}