package utils

import (
	"bank-backend/module/bank/entity"
	"time"

	"github.com/google/uuid"
)

func TopUpDTO(user entity.User, prev int, tid uuid.UUID, topup int, time time.Time) entity.TopUpResponse {
	response := entity.TopUpResponse{
		TopUpId:       tid.String(),
		BalanceBefore: prev,
		BalanceAfter:  user.Balance,
		Amount:        topup,
		CreatedAt:     time.String(),
	}
	return response
}

func PaymentDTO(user entity.User, prev int, tid uuid.UUID, topup int, time time.Time, remarks string) entity.PaymentResponse {
	response := entity.PaymentResponse{
		PaymentID:     tid.String(),
		BalanceBefore: prev,
		BalanceAfter:  user.Balance,
		Remarks:       remarks,
		Amount:        topup,
		CreatedAt:     time.String(),
	}
	return response
}

func TransferDTO(balanceAfter int, prev int, tid uuid.UUID, topup int, time string, remarks string, targetTransfer string) entity.TransferResponse {
	response := entity.TransferResponse{
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
