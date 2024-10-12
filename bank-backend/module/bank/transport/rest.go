package transport

import (
	"bank-backend/module/bank/config"
	"bank-backend/module/bank/entity"
	"bank-backend/module/bank/internal/queue"
	"bank-backend/module/bank/internal/repository"
	"bank-backend/module/bank/internal/usecase"
	"bank-backend/module/middleware"
	"bank-backend/pkg"
	"bank-backend/utils"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type Rest struct {
	bankUC   usecase.BankUseCase
	validate *validator.Validate
}

func NewRest(cfg config.BankConfig) {
	bankRepo := repository.NewBankRepository(cfg.PGx)
	processTransferQueue := queue.NewProcessTransferQueue(*cfg.Producer, cfg.ProcessTranferTopic)
	bankUsecase := usecase.NewBankUseCase(*bankRepo, processTransferQueue)
	transport := Rest{bankUC: bankUsecase, validate: cfg.Validate}

	// Initialize Fiber app
	transport.mountBank(cfg.Fiber)

}
func (r *Rest) mountBank(app *fiber.App) {

	fmt.Println("bank")
	fmt.Println(app)
	app.Post("/api/v1/topup", r.Topup, middleware.JwtMiddleware(), middleware.RoleBasedMiddleware())
	app.Post("/api/v1/payment", r.Payment, middleware.JwtMiddleware(), middleware.RoleBasedMiddleware())
	app.Post("/api/v1/transfer", r.Transfer, middleware.JwtMiddleware(), middleware.RoleBasedMiddleware())
}

func (r *Rest) Topup(ctx fiber.Ctx) error {

	var (
		lvState1       = utils.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = utils.LogEventStateCallUsecase
		lfState2Status = "state_2_call_usecase"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	// Retrieve the user phoneNumber from the context
	userPhoneNumber := ctx.Locals("user-phone").(string)

	lf = append(lf, pkg.LogEventState(lvState1))
	topupPayload := new(entity.TopUpRequest)
	err := ctx.Bind().JSON(topupPayload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = r.validate.Struct(topupPayload); err != nil {
		errors := utils.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(utils.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(topupPayload),
	)

	res, err := r.bankUC.Topup(ctx, *topupPayload, userPhoneNumber)
	lf = append(lf, pkg.LogEventState(lvState2))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(utils.StandardResponse{
		Status: "SUCCESS",
		Result: res,
	})
}

func (r *Rest) Payment(ctx fiber.Ctx) error {

	var (
		lvState1       = utils.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = utils.LogEventStateCallUsecase
		lfState2Status = "state_2_call_usecase"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	// Retrieve the user phoneNumber from the context
	userPhoneNumber := ctx.Locals("user-phone").(string)

	lf = append(lf, pkg.LogEventState(lvState1))
	payment := new(entity.PaymentRequest)
	err := ctx.Bind().JSON(payment)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = r.validate.Struct(payment); err != nil {
		errors := utils.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(utils.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(payment),
	)

	res, err := r.bankUC.Payment(ctx, *payment, userPhoneNumber)
	lf = append(lf, pkg.LogEventState(lvState2))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(utils.StandardResponse{
		Status: "SUCCESS",
		Result: res,
	})
}

func (r *Rest) Transfer(ctx fiber.Ctx) error {

	var (
		lvState1       = utils.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = utils.LogEventStateCallUsecase
		lfState2Status = "state_2_call_usecase"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	// Retrieve the user phoneNumber from the context
	userPhoneNumber := ctx.Locals("user-phone").(string)

	lf = append(lf, pkg.LogEventState(lvState1))
	transferPayload := new(entity.TransferRequest)
	err := ctx.Bind().JSON(transferPayload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = r.validate.Struct(transferPayload); err != nil {
		errors := utils.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(utils.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(transferPayload),
	)
	res, err := r.bankUC.Transfer(ctx, *transferPayload, userPhoneNumber)
	lf = append(lf, pkg.LogEventState(lvState2))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(utils.StandardResponse{
		Status: "SUCCESS",
		Result: res,
	})
}
