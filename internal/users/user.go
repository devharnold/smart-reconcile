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
