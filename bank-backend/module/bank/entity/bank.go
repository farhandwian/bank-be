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

type TopUpRequest struct {
	Amount int `json:"amount"validate:"required,min=1,numeric"`
}

type TopUpResponse struct {
	TopUpId       string `json:"top_up_id"`
	BalanceBefore int    `json:"balance_before"`
	BalanceAfter  int    `json:"balance_after"`
	Amount        int    `json:"amount"`
	Remarks       string `json:"remarks,omitempty"`
	CreatedAt     string `json:"created_at"`
}

type PaymentRequest struct {
	Amount  int    `json:"amount"validate:"required,min=1,numeric"`
	Remarks string `json:"remarks" validate:"required,max=50"`
}

type PaymentResponse struct {
	PaymentID     string `json:"payment_id"`
	BalanceBefore int    `json:"balance_before"`
	BalanceAfter  int    `json:"balance_after"`
	Amount        int    `json:"amount"`
	Remarks       string `json:"remarks,omitempty"`
	CreatedAt     string `json:"created_at"`
}

type TransferRequest struct {
	Amount     int    `json:"amount"validate:"required,min=1,numeric"`
	TargetUser string `json:"target_user"validate:"required"`
	Remarks    string `json:"remarks" validate:"required,max=50"`
}

type TransferResponse struct {
	TransferID     string `json:"transfer_id"`
	BalanceBefore  int    `json:"balance_before"`
	BalanceAfter   int    `json:"balance_after"`
	TargetTransfer string `json:"target_transfer"`
	Amount         int    `json:"amount"`
	Remarks        string `json:"remarks,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type TransferEvent struct {
	Transfer              string `json:"transaction_id"`
	Amount                int    `json:"amount"`
	PhoneNumberOriginUser string `json:"phone_number_origin_user"`
	TargetUser            string `json:"target_user"`
	Remarks               string `json:"remarks"`
	CreatedAt             string `json:"created_at"`
}
