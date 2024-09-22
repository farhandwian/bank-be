package user

type RegisterRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=1,max=20,alphanum"`
	LastName    string `json:"last_name" validate:"required,min=1,max=20,alphanum"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,indonesianphone"`
	Pin         string `json:"pin" validate:"required,len=6,numeric"`
}

type UpdateRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=20,alphanum"`
	LastName  string `json:"last_name" validate:"required,min=1,max=20,alphanum"`
	Address   string `json:"address" validate:"required"`
}

type RegisterResponse struct {
	UserID      string `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   string `json:"created_at"`
}
type UpdateResponse struct {
	UserID      string `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Updated_at  string `json:"updated_at"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"validate:"required,indonesianphone"`
	Pin         string `json:"pin"validate:"required,len=6,numeric"`
}

type LoginResponse struct {
	Token        string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
