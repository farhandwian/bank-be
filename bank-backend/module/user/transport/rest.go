package user

import (
	"bank-backend/module/middleware"
	"bank-backend/module/user/config"
	"bank-backend/module/user/entity"
	"bank-backend/module/user/internal/repository"
	"bank-backend/module/user/internal/usecase"
	"errors"
	"fmt"

	"bank-backend/pkg"
	"bank-backend/utils"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type Rest struct {
	userUC   usecase.UserUsecase
	validate *validator.Validate
}

func NewRest(cfg config.UserConfig) {
	userRepo := repository.NewUserRepository(cfg.PGx)
	userUsecase := usecase.NewUserUseCase(*userRepo)
	transport := Rest{userUC: userUsecase, validate: cfg.Validate}
	// Initialize Fiber app
	transport.mountUser(cfg.Fiber)

}

func (r *Rest) mountUser(app *fiber.App) {

	fmt.Println("user")
	fmt.Println(app)
	app.Post("/api/v1/register", r.Register)
	app.Post("/api/v1/login", r.Login)
	app.Post("/api/v1/refresh", r.RefreshToken)
	app.Put("/api/v1/update", r.UpdateProfile, middleware.JwtMiddleware(), middleware.RoleBasedMiddleware())
}

func (r *Rest) Register(ctx fiber.Ctx) error {

	var (
		lvState1       = utils.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = utils.LogEventStateCallUsecase
		lfState2Status = "state_2_call_usecase"
		lf             = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))
	registerPayload := new(entity.RegisterRequest)
	err := ctx.Bind().JSON(registerPayload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = r.validate.Struct(registerPayload); err != nil {
		errors := utils.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(utils.StandardResponse{Errors: errors})
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(registerPayload),
	)

	res, err := r.userUC.Register(ctx, *registerPayload)
	lf = append(lf, pkg.LogEventState(lvState2))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(utils.StandardResponse{
		Status: "SUCCESS",
		Result: res,
	})
}

func (r *Rest) Login(ctx fiber.Ctx) error {

	var (
		lvState1       = utils.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = utils.LogEventStateCallUsecase
		lfState2Status = "state_2_call_usecase"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))
	loginPayload := new(entity.LoginRequest)
	err := ctx.Bind().JSON(loginPayload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = r.validate.Struct(loginPayload); err != nil {
		errors := utils.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(utils.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(loginPayload),
	)
	res, err := r.userUC.Login(ctx, *loginPayload)
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

func (r *Rest) RefreshToken(ctx fiber.Ctx) error {
	var (
		lvState1       = utils.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = utils.LogEventStateCallUsecase
		lfState2Status = "state_2_call_usecase"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))
	refreshPayload := new(entity.RefreshRequest)
	err := ctx.Bind().JSON(refreshPayload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: "error processed request",
		})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(refreshPayload),
	)

	res, err := r.userUC.RefreshToken(ctx, *refreshPayload)
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

func (r *Rest) UpdateProfile(ctx fiber.Ctx) error {

	var (
		lvState1       = utils.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	// Retrieve the user phoneNumber from the context
	userPhoneNumber, ok := ctx.Locals("user-phone").(string)
	if !ok {
		// Handle the case where the value is not a string
		return errors.New("user-phone is not a string")
	}

	lf = append(lf, pkg.LogEventState(lvState1))
	updatePayload := new(entity.UpdateProfileRequest)
	err := ctx.Bind().JSON(updatePayload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(utils.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = r.validate.Struct(updatePayload); err != nil {
		errors := utils.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(utils.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(updatePayload),
	)

	res, err := r.userUC.UpdateProfile(ctx, *updatePayload, userPhoneNumber)
	return ctx.Status(http.StatusOK).JSON(utils.StandardResponse{
		Status: "SUCCESS",
		Result: res,
	})
}
