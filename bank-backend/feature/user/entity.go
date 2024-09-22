package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FirstName   string
	LastName    string
	Version     int
	PhoneNumber string
	Address     string
	Pin         string
}
