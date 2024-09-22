package bank

import (
	"github.com/google/uuid"
	"time"
)

func topUpDTO(user User, prev int, tid uuid.UUID, topup int, time time.Time) TopUpResponse {
	response := TopUpResponse{
		TopUpId:       tid.String(),
		BalanceBefore: prev,
		BalanceAfter:  user.Balance,
		Amount:        topup,
		CreatedAt:     time.String(),
	}
	return response
}

func paymentDTO(user User, prev int, tid uuid.UUID, topup int, time time.Time, remarks string) PaymentResponse {
	response := PaymentResponse{
		PaymentID:     tid.String(),
		BalanceBefore: prev,
		BalanceAfter:  user.Balance,
		Remarks:       remarks,
		Amount:        topup,
		CreatedAt:     time.String(),
	}
	return response
}

func transferDTO(balanceAfter int, prev int, tid uuid.UUID, topup int, time string, remarks string, targetTransfer string) TransferResponse {
	response := TransferResponse{
		TransferID:     tid.String(),
		BalanceBefore:  prev,
		BalanceAfter:   balanceAfter,
		Remarks:        remarks,
		TargetTransfer: targetTransfer,
		Amount:         topup,
		CreatedAt:      time,
	}
	return response
}
