package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FirstName   string
	LastName    string
	Version     int
	PhoneNumber string
	Balance     int
	Address     string
	Pin         string
}

type RegisterRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=1,max=20,alphanum"`
	LastName    string `json:"last_name" validate:"required,min=1,max=20,alphanum"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,indonesianphone"`
	Pin         string `json:"pin" validate:"required,len=6,numeric"`
}

type UpdateProfileRequest struct {
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

type Meta struct {
	HTTPStatus int `json:"http_status"`
}

// type RegisterResponse struct {
// 	Message  string               `json:"message"`
// 	Meta     Meta                 `json:"meta"`
// 	Register RegisterResponse `json:"register,omitempty"`
// }

// bisa dilakukan seperti ini jika banyak perbedaan response antar entity
// func NewRegisterResponse(message string, code int, register RegisterResponse) *RegisterResponse {
// 	res := &RegisterResponse{
// 		Message:  message,
// 		Meta:     Meta{HTTPStatus: code},
// 		Register: register,
// 	}

// 	return res
// }

type UpdateProfileResponse struct {
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
