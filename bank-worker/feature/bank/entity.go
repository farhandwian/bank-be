package bank

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
	Balance     int
	Address     string
	Pin         string
}

type Transaction struct {
	ID              uuid.UUID
	Remarks         string
	Amount          int
	BalanceBefore   int
	BalanceAfter    int
	TransactionType string
	UserID          uuid.UUID
	TargetUserID    uuid.UUID
	CreatedDate     time.Time
	Version         int
}
