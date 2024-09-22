package bank

import (
	"bank-backend/feature/middleware"
	"bank-backend/feature/shared"
	"bank-backend/pkg"
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

func FiberRoute(app *fiber.App) {
	app.Post("/api/v1/topup", Topup, middleware.JwtMiddleware(), middleware.RoleBasedMiddleware())
	app.Post("/api/v1/payment", Payment, middleware.JwtMiddleware(), middleware.RoleBasedMiddleware())
	app.Post("/api/v1/transfer", Transfer, middleware.JwtMiddleware(), middleware.RoleBasedMiddleware())
}

func Topup(ctx fiber.Ctx) error {

	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateUpdateDB
		lfState2Status = "state_2_update_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	// Retrieve the user phoneNumber from the context
	userPhoneNumber := ctx.Locals("user-phone")

	lf = append(lf, pkg.LogEventState(lvState1))
	topup := new(TopUpRequest)
	err := ctx.Bind().JSON(topup)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = validate.Struct(topup); err != nil {
		errors := shared.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(topup),
	)
	/*------------------------------------
	| Step 2 : Update Balance
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	u := User{
		UpdatedAt:   time.Now(),
		Balance:     topup.Amount,
		PhoneNumber: userPhoneNumber.(string),
	}

	user, prev, tid, createdAt, err := updateTopUpt(ctx.Context(), u)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "user not found",
		})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(user),
	)

	dto := topUpDTO(user, prev, tid, topup.Amount, createdAt)

	return ctx.Status(http.StatusOK).JSON(shared.StandardResponse{
		Status: "SUCCESS",
		Result: dto,
	})
}

func Payment(ctx fiber.Ctx) error {

	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateUpdateDB
		lfState2Status = "state_2_update_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	// Retrieve the user phoneNumber from the context
	userPhoneNumber := ctx.Locals("user-phone")

	lf = append(lf, pkg.LogEventState(lvState1))
	payment := new(PaymentRequest)
	err := ctx.Bind().JSON(payment)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = validate.Struct(payment); err != nil {
		errors := shared.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(payment),
	)
	/*------------------------------------
	| Step 2 : Update Balance
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	u := User{
		UpdatedAt:   time.Now(),
		Balance:     payment.Amount,
		PhoneNumber: userPhoneNumber.(string),
	}

	user, prev, tid, createdAt, err := updatePayment(ctx.Context(), u, payment.Remarks)
	if err != nil {
		if err == errBalanceNotEnough {
			lf = append(lf, pkg.LogStatusFailed(lfState2Status))
			pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
			return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{
				Message: "Balance is not enough",
			})
		}
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "user not found",
		})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(user),
	)

	dto := paymentDTO(user, prev, tid, payment.Amount, createdAt, payment.Remarks)

	return ctx.Status(http.StatusOK).JSON(shared.StandardResponse{
		Status: "SUCCESS",
		Result: dto,
	})
}

func Transfer(ctx fiber.Ctx) error {

	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateFetchDB
		lfState2Status = "state_2_fetch_db_status"

		lvState3       = shared.LogEventStateKafkaPublish
		lfState3Status = "state_2_kafka_publish_status"

		lf = []slog.Attr{
			pkg.LogEventName("bank-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	// Retrieve the user phoneNumber from the context
	userPhoneNumber := ctx.Locals("user-phone")
	strPhoneNumber := userPhoneNumber.(string)

	lf = append(lf, pkg.LogEventState(lvState1))
	transfer := new(TransferRequest)
	err := ctx.Bind().JSON(transfer)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = validate.Struct(transfer); err != nil {
		errors := shared.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(transfer),
	)
	/*------------------------------------
	| Step 2 : Fetch User and Check Balance
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	// check origin user
	originUser, err := checkIfUserExistByPhoneNumber(ctx.Context(), strPhoneNumber)
	if err != nil {
		if err == errUserNotFound {
			lf = append(lf, pkg.LogStatusFailed(lfState2Status))
			pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
			return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
				Message: "user not found",
			})
		}
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}

	if originUser.Balance < transfer.Amount {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "balance is not enough", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "user balance not enough",
		})
	}

	//check destination user
	parse, err := uuid.Parse(transfer.TargetUser)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "Internal server error",
		})
	}
	_, err = checkIfUserExistByID(ctx.Context(), parse)
	if err != nil {
		if err == errUserNotFound {
			lf = append(lf, pkg.LogStatusFailed(lfState2Status))
			pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
			return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
				Message: "target user not found",
			})
		}
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(originUser),
	)

	/*------------------------------------
	| Step 3 : Publish TransferEvent
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	id, err := pkg.GenerateId()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), "generate uuid error", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}
	now := time.Now()
	// Format the time
	formatted := now.Format("2006-01-02 15:04:05.000000")

	event := TransferEvent{
		Transfer:              id.String(),
		Amount:                transfer.Amount,
		PhoneNumberOriginUser: userPhoneNumber.(string),
		TargetUser:            transfer.TargetUser,
		Remarks:               transfer.Remarks,
		CreatedAt:             formatted,
	}
	messageByte, err := json.Marshal(event)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), "kafka publish error", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}
	err = pkg.PublishMessage(kp, CreateNewTransferTopic, string(messageByte))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), "kafka publish error", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}

	dto := transferDTO(originUser.Balance-transfer.Amount, originUser.Balance, id, transfer.Amount, event.CreatedAt, transfer.Remarks, transfer.TargetUser)

	return ctx.Status(http.StatusOK).JSON(shared.StandardResponse{
		Status: "SUCCESS",
		Result: dto,
	})
}
