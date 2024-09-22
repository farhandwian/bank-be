package user

import (
	"bank-backend/feature/middleware"
	"bank-backend/feature/shared"
	"bank-backend/pkg"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"time"
)

func FiberRoute(app *fiber.App) {
	app.Post("/api/v1/register", Register)
	app.Post("/api/v1/login", Login)
	app.Post("/api/v1/refresh", RefreshToken)
	app.Put("/api/v1/update", UpdateProfile, middleware.JwtMiddleware(), middleware.RoleBasedMiddleware())
}

func Register(ctx fiber.Ctx) error {

	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateFetchDB
		lfState2Status = "state_2_fetch_user_db_status"

		lvState3       = shared.LogEventStateInsertDB
		lfState3Status = "state_3_insert_user_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))
	register := new(RegisterRequest)
	err := ctx.Bind().JSON(register)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = validate.Struct(register); err != nil {
		errors := shared.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{Errors: errors})
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(register),
	)
	/*------------------------------------
	| Step 2 : Check If Username Is Exist
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	exists, _, _, err := checkPhoneNumberExists(ctx.Context(), register.PhoneNumber)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}

	if exists {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "phone number already registered", err, lf)
		return ctx.Status(http.StatusNotFound).JSON(shared.StandardResponse{
			Message: "Phone Number already registered",
		})
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(register),
	)

	/*------------------------------------
	| Step 3 : Insert Register User
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	password, err := bcrypt.GenerateFromPassword([]byte(register.Pin), bcrypt.DefaultCost)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}

	id, err := pkg.GenerateId()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}
	user := User{
		ID:          id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     1,
		PhoneNumber: register.PhoneNumber,
		FirstName:   register.FirstName,
		LastName:    register.LastName,
		Address:     register.Address,
		Balance:     0,
		Pin:         string(password),
	}

	rUser, err := insertUser(ctx.Context(), user)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: err.Error(),
		})
	}

	dto := userToDTO(rUser)

	lf = append(lf,
		pkg.LogStatusSuccess(lfState3Status),
		pkg.LogEventPayload(register),
	)

	return ctx.Status(http.StatusCreated).JSON(shared.StandardResponse{
		Status: "SUCCESS",
		Result: dto,
	})
}

func Login(ctx fiber.Ctx) error {

	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateFetchDB
		lfState2Status = "state_2_fetch_user_db_status"

		lvState3       = shared.LogEventStateFetchDB
		lfState3Status = "state_3_compare_password_status"

		lvState4       = shared.LogEventStateSetToken
		lfState4Status = "state_4_set_token_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))
	login := new(LoginRequest)
	err := ctx.Bind().JSON(login)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = validate.Struct(login); err != nil {
		errors := shared.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(login),
	)
	/*------------------------------------
	| Step 2 : Check If Username Is Exist
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))
	exists, PhoneNumber, password, err := checkPhoneNumberExists(ctx.Context(), login.PhoneNumber)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "internal server error",
		})
	}

	if !exists {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "Phone Number and PIN doesn't match", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{
			Message: "Phone Number and PIN doesn't match",
		})
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload("success fetch data"),
	)

	/*------------------------------------
	| Step 3 : Password
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(login.Pin)); err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), "Phone Number and PIN doesn't match", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{
			Message: "Phone Number and PIN doesn't match",
		})
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState3Status),
		pkg.LogEventPayload(PhoneNumber),
	)

	/*------------------------------------
	| Step 4 : Generate Access Token & refresh Token
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState4))
	accessToken, err := pkg.GenerateAccessTokens(PhoneNumber)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState4Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "Failed to generate refresh tokens",
		})
	}

	refreshToken, err := pkg.GenerateRefreshTokens(PhoneNumber)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lvState4))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "Failed to generate refresh tokens",
		})
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState4Status),
		pkg.LogEventPayload(refreshToken),
	)

	return ctx.Status(http.StatusOK).JSON(shared.StandardResponse{
		Status: "SUCCESS",
		Result: LoginResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
		},
	})
}

func RefreshToken(ctx fiber.Ctx) error {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateValidateToken
		lfState2Status = "state_2_validated_token_status"

		lvState3       = shared.LogEventStateSetToken
		lfState3Status = "state_3_set_token_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))
	refresh := new(RefreshRequest)
	err := ctx.Bind().JSON(refresh)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "error processed request",
		})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(refresh),
	)

	/*------------------------------------
	| Step 2 : Validate Token
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	// Parse and validate the refresh token
	token, err := jwt.ParseWithClaims(refresh.RefreshToken, &pkg.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(pkg.JWTSecret), nil
	})

	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusUnauthorized).JSON(shared.StandardResponse{
			Message: "Invalid refresh token",
		})
	}

	claims, ok := token.Claims.(*pkg.Claims)
	if !ok || !token.Valid {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "invalid refresh token", errors.New("invalid refresh token"), lf)
		return ctx.Status(http.StatusUnauthorized).JSON(shared.StandardResponse{
			Message: "Invalid refresh token claims",
		})
	}

	// Check if the token is expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "Refresh token has expired", errors.New("Refresh token has expired"), lf)
		return ctx.Status(http.StatusUnauthorized).JSON(shared.StandardResponse{
			Message: "Refresh token has expired",
		})
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(claims),
	)
	/*------------------------------------
	| Step 3 : Reg-generate Access Token
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	accessToken, err := pkg.GenerateAccessTokens(claims.PhoneNumber)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return ctx.Status(http.StatusUnauthorized).JSON(shared.StandardResponse{
			Message: "Failed to generate access tokens",
		})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState3Status),
		pkg.LogEventPayload(accessToken),
	)

	return ctx.Status(http.StatusOK).JSON(shared.StandardResponse{
		Status: "SUCCESS",
		Result: LoginResponse{Token: accessToken},
	})
}

func UpdateProfile(ctx fiber.Ctx) error {

	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_request_status"

		lvState2       = shared.LogEventStateUpdateDB
		lfState2Status = "state_2_update_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	// Retrieve the user phoneNumber from the context
	userPhoneNumber := ctx.Locals("user-phone")

	lf = append(lf, pkg.LogEventState(lvState1))
	update := new(UpdateRequest)
	err := ctx.Bind().JSON(update)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "error processed request", err, lf)
		return ctx.Status(http.StatusInternalServerError).JSON(shared.StandardResponse{
			Message: "error processed request",
		})
	}
	// Validate the struct
	if err = validate.Struct(update); err != nil {
		errors := shared.FormatValidationErrors(err)
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx.Context(), "validation invalid", err, lf)
		return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{Errors: errors})
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(update),
	)
	/*------------------------------------
	| Step 2 : Update User
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	u := User{
		UpdatedAt:   time.Now(),
		FirstName:   update.FirstName,
		LastName:    update.LastName,
		PhoneNumber: userPhoneNumber.(string),
		Address:     update.Address,
	}
	user, err := updateUser(ctx.Context(), u)
	if err != nil {
		if err == errUserNotFound {
			lf = append(lf, pkg.LogStatusFailed(lfState2Status))
			pkg.LogWarnWithContext(ctx.Context(), "user not found", err, lf)
			return ctx.Status(http.StatusBadRequest).JSON(shared.StandardResponse{
				Message: "user not found",
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
		pkg.LogEventPayload(update),
	)
	dto := userUpdataeToDTO(user)

	return ctx.Status(http.StatusOK).JSON(shared.StandardResponse{
		Status: "SUCCESS",
		Result: dto,
	})
}
