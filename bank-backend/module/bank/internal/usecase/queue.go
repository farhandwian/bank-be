package usecase

import (
	"bank-backend/module/bank/entity"
	"context"

	"github.com/google/uuid"
)

type ProcessTransferQueue interface {
	PublishProcessTransferJob(ctx context.Context, request entity.TransferRequest, userPhoneNumber string) (uuid.UUID, string, error)
}
