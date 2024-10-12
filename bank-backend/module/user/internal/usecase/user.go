package usecase

import (
	"bank-backend/module/user/entity"
	"bank-backend/module/user/internal/repository"
	"bank-backend/module/user/utils"
	"bank-backend/pkg"
	utls "bank-backend/utils"
	"bank-backend/utils/pgsql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(ctx fiber.Ctx, request entity.RegisterRequest) (entity.RegisterResponse, error)
	Login(ctx fiber.Ctx, request entity.LoginRequest) (entity.LoginResponse, error)
	RefreshToken(ctx fiber.Ctx, request entity.RefreshRequest) (entity.LoginResponse, error)
	UpdateProfile(ctx fiber.Ctx, request entity.UpdateProfileRequest, userPhoneNumber string) (entity.UpdateProfileResponse, error)
}

type UserUC struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUC {
	return &UserUC{userRepo: userRepo}
}

func (u *UserUC) Register(ctx fiber.Ctx, request entity.RegisterRequest) (entity.RegisterResponse, error) {
	var (
		lvState2       = utls.LogEventStateFetchDB
		lfState2Status = "state_2_fetch_user_db_status"

		lvState3       = utls.LogEventStateInsertDB
		lfState3Status = "state_3_insert_user_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 2 : Check If Username Is Exist
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	exists, _, _, err := u.userRepo.CheckPhoneNumberExists(ctx.Context(), request.PhoneNumber)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		return entity.RegisterResponse{}, err
	}

	if exists {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "phone number already registered", err, lf)
		return entity.RegisterResponse{}, err
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(request),
	)

	/*------------------------------------
	| Step 3 : Insert Register User
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	password, err := bcrypt.GenerateFromPassword([]byte(request.Pin), bcrypt.DefaultCost)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.RegisterResponse{}, err
	}

	id, err := pkg.GenerateId()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.RegisterResponse{}, err
	}
	user := entity.User{
		ID:          id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     1,
		PhoneNumber: request.PhoneNumber,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Address:     request.Address,
		Balance:     0,
		Pin:         string(password),
	}

	rUser, err := u.userRepo.InsertUser(ctx.Context(), user)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.RegisterResponse{}, err
	}

	dto := utils.UserToDTO(rUser)

	lf = append(lf,
		pkg.LogStatusSuccess(lfState3Status),
		pkg.LogEventPayload(request),
	)

	return dto, nil
}

func (u *UserUC) Login(ctx fiber.Ctx, request entity.LoginRequest) (entity.LoginResponse, error) {
	var (
		lvState2       = utls.LogEventStateFetchDB
		lfState2Status = "state_2_fetch_user_db_status"

		lvState3       = utls.LogEventStateFetchDB
		lfState3Status = "state_3_compare_password_status"

		lvState4       = utls.LogEventStateSetToken
		lfState4Status = "state_4_set_token_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)

	/*------------------------------------
	| Step 2 : Check If Username Is Exist
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))
	exists, PhoneNumber, password, err := u.userRepo.CheckPhoneNumberExists(ctx.Context(), request.PhoneNumber)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.LoginResponse{}, err
	}

	if !exists {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "Phone Number and PIN doesn't match", err, lf)
		return entity.LoginResponse{}, err
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload("success fetch data"),
	)

	/*------------------------------------
	| Step 3 : Password
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(request.Pin)); err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogWarnWithContext(ctx.Context(), "Phone Number and PIN doesn't match", err, lf)
		return entity.LoginResponse{}, err
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
		return entity.LoginResponse{}, nil
	}

	refreshToken, err := pkg.GenerateRefreshTokens(PhoneNumber)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lvState4))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.LoginResponse{}, err
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState4Status),
		pkg.LogEventPayload(refreshToken),
	)
	res := entity.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	return res, nil
}

func (u *UserUC) RefreshToken(ctx fiber.Ctx, request entity.RefreshRequest) (entity.LoginResponse, error) {
	var (
		lvState2       = utls.LogEventStateValidateToken
		lfState2Status = "state_2_validated_token_status"

		lvState3       = utls.LogEventStateSetToken
		lfState3Status = "state_3_set_token_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)

	/*------------------------------------
	| Step 2 : Validate Token
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	// Parse and validate the refresh token
	token, err := jwt.ParseWithClaims(request.RefreshToken, &pkg.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(pkg.JWTSecret), nil
	})

	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.LoginResponse{}, err
	}

	claims, ok := token.Claims.(*pkg.Claims)
	if !ok || !token.Valid {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "invalid refresh token", errors.New("invalid refresh token"), lf)
		return entity.LoginResponse{}, err
	}

	// Check if the token is expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), "Refresh token has expired", errors.New("Refresh token has expired"), lf)
		return entity.LoginResponse{}, err
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
		return entity.LoginResponse{}, err
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState3Status),
		pkg.LogEventPayload(accessToken),
	)
	res := entity.LoginResponse{
		Token: accessToken,
	}

	return res, nil
}

func (u *UserUC) UpdateProfile(ctx fiber.Ctx, request entity.UpdateProfileRequest, userPhoneNumber string) (entity.UpdateProfileResponse, error) {
	var (
		lvState2       = utls.LogEventStateUpdateDB
		lfState2Status = "state_2_update_db_status"

		lf = []slog.Attr{
			pkg.LogEventName("user-service"),
		}
	)
	/*------------------------------------
	| Step 2 : Update User
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	userEntity := entity.User{
		UpdatedAt:   time.Now(),
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: userPhoneNumber,
		Address:     request.Address,
	}
	user, err := u.userRepo.UpdateUser(ctx.Context(), userEntity)
	if err != nil {
		if err == pgsql.ErrUserNotFound {
			lf = append(lf, pkg.LogStatusFailed(lfState2Status))
			pkg.LogWarnWithContext(ctx.Context(), "user not found", err, lf)
			return entity.UpdateProfileResponse{}, err
		}
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogWarnWithContext(ctx.Context(), err.Error(), err, lf)
		return entity.UpdateProfileResponse{}, err
	}
	lf = append(lf,
		pkg.LogStatusSuccess(lfState2Status),
		pkg.LogEventPayload(user),
	)
	dto := utils.UserUpdateToDTO(user)
	return dto, nil
}
