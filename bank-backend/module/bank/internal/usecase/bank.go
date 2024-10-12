package usecase

import (
	"bank-backend/module/bank/entity"
	"bank-backend/module/bank/internal/repository"
	"bank-backend/module/bank/utils"
	"bank-backend/pkg"
	utls "bank-backend/utils"
	"bank-backend/utils/pgsql"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type BankUseCase interface {
	Topup(ctx fiber.Ctx, request entity.TopUpRequest, userPhoneNumber string) (entity.TopUpResponse, error)
	Payment(ctx fiber.Ctx, request entity.PaymentRequest, userPhoneNumber string) (entity.PaymentResponse, error)
	Transfer(ctx fiber.Ctx, request entity.TransferRequest, userPhoneNumber string) (entity.TransferResponse, error)
}

type BankUC struct {
	bankRepo        repository.BankRepository
	processTransfer ProcessTransferQueue
}

func NewBankUseCase(bankRepo repository.BankRepository, processTransfer ProcessTransferQueue) *BankUC {
	return &BankUC{bankRepo: bankRepo, processTransfer: processTransfer}
}

func (b *BankUC) Topup(ctx fiber.Ctx, request entity.TopUpRequest, userPhoneNumber string) (entity.TopUpResponse, error) {
	var (
		lvState2       = utls.LogEventStateUpdateDB
		lfState2Status = "state_2_update_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 2 : Update Balance
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	u := entity.User{
		UpdatedAt:   time.Now(),
		Balance:     request.Amount,
		PhoneNumber: userPhoneNumber,
	}

	user, prev, tid, createdAt, err := b.bankRepo.UpdateTopUpt(ctx.Context(), u)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.TopUpResponse{}, err
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(user),
	)

	dto := utils.TopUpDTO(user, prev, tid, request.Amount, createdAt)

	return dto, nil
}

func (b *BankUC) Payment(ctx fiber.Ctx, request entity.PaymentRequest, userPhoneNumber string) (entity.PaymentResponse, error) {

	var (
		lvState2       = utls.LogEventStateUpdateDB
		lfState2Status = "state_2_update_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 2 : Update Balance
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	u := entity.User{
		UpdatedAt:   time.Now(),
		Balance:     request.Amount,
		PhoneNumber: userPhoneNumber,
	}

	user, prev, tid, createdAt, err := b.bankRepo.UpdatePayment(ctx.Context(), u, request.Remarks)
	if err != nil {
		if err == pgsql.ErrBalanceNotEnough {
			lf = append(lf, pkg.LogStatusFailed(lfState2Status))
			pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
			return entity.PaymentResponse{}, err
		}
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.PaymentResponse{}, err
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(user),
	)

	dto := utils.PaymentDTO(user, prev, tid, request.Amount, createdAt, request.Remarks)

	return dto, nil
}

func (b *BankUC) Transfer(ctx fiber.Ctx, request entity.TransferRequest, userPhoneNumber string) (entity.TransferResponse, error) {
	var (
		lvState2       = utls.LogEventStateFetchDB
		lfState2Status = "state_2_fetch_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)

	/*------------------------------------
	| Step 2 : Fetch entity.User and Check Balance
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	// check origin user
	originUser, err := b.bankRepo.CheckIfUserExistByPhoneNumber(ctx.Context(), userPhoneNumber)
	if err != nil {
		if err == pgsql.ErrUserNotFound {
			lf = append(lf, pkg.LogStatusFailed(lfState2Status))
			pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
			return entity.TransferResponse{}, err
		}
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.TransferResponse{}, err
	}

	if originUser.Balance < request.Amount {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "balance is not enough", err, lf)
		return entity.TransferResponse{}, err
	}

	//check destination user
	parse, err := uuid.Parse(request.TargetUser)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.TransferResponse{}, err
	}
	_, err = b.bankRepo.CheckIfUserExistByID(ctx.Context(), parse)
	if err != nil {
		if err == pgsql.ErrUserNotFound {
			lf = append(lf, pkg.LogStatusFailed(lfState2Status))
			pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
			return entity.TransferResponse{}, err
		}
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.TransferResponse{}, err
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(originUser),
	)

	/*------------------------------------
	| Step 3 : Publish TransferEvent
	* ----------------------------------*/
	id, created_at, err := b.processTransfer.PublishProcessTransferJob(ctx.Context(), request, userPhoneNumber)

	dto := utils.TransferDTO(originUser.Balance-request.Amount, originUser.Balance, id, request.Amount, created_at, request.Remarks, request.TargetUser)

	return dto, nil
}
