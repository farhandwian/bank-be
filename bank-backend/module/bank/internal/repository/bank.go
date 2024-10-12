package repository

import (
	"context"
	"time"

	"bank-backend/module/bank/entity"
	"bank-backend/pkg"
	"bank-backend/utils/pgsql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"
)

type BankRepository struct {
	db *pgxpool.Pool
}

func NewBankRepository(db *pgxpool.Pool) *BankRepository {
	return &BankRepository{db: db}
}

func (b *BankRepository) CheckIfUserExistByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error) {
	user := entity.User{}
	query := `SELECT id, balance FROM "user" where phone_number = $1`

	err := b.db.QueryRow(ctx, query, phoneNumber).Scan(&user.ID, &user.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = pgsql.ErrUserNotFound
			return user, err
		}
		return user, err
	}
	return user, nil
}

func (b *BankRepository) CheckIfUserExistByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	user := entity.User{}
	query := `SELECT id, balance FROM "user" where id = $1`

	err := b.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = pgsql.ErrUserNotFound
			return user, err
		}
		return user, err
	}
	return user, nil
}

func (b *BankRepository) UpdateTopUpt(ctx context.Context, user entity.User) (entity.User, int, uuid.UUID, time.Time, error) {
	returningUser := entity.User{}
	tx, err := b.db.Begin(ctx)
	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	defer tx.Rollback(ctx)

	query := `update "user" set balance = balance + $1, version = version+1, updated_at = $2 where phone_number = $3 and version = $4 RETURNING id, balance, updated_at `

	selectUser := `select phone_number, balance, version from "user" where phone_number = $1`

	transactionQuery := `
		INSERT INTO transaction (id, amount, balance_before, balance_after, transaction_type, user_id, created_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at
	`

	var PhoneNumber string
	var prevBalance int
	var Version int

	err = tx.QueryRow(ctx, selectUser, user.PhoneNumber).Scan(&PhoneNumber, &prevBalance, &Version)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = pgsql.ErrUserNotFound
		}
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	err = tx.QueryRow(ctx, query, user.Balance, time.Now(), PhoneNumber, Version).Scan(&returningUser.ID, &returningUser.Balance, &returningUser.UpdatedAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	id, err := pkg.GenerateId()
	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	// Sample transaction data
	transaction := entity.Transaction{
		ID:              id,
		Amount:          user.Balance,
		BalanceBefore:   prevBalance,
		BalanceAfter:    returningUser.Balance,
		TransactionType: "CREDIT",
		UserID:          returningUser.ID,
		CreatedDate:     time.Now(),
		Version:         1,
	}

	var transactionId uuid.UUID
	var createdAt time.Time

	err = tx.QueryRow(context.Background(), transactionQuery,
		transaction.ID,
		transaction.Amount,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.TransactionType,
		transaction.UserID,
		transaction.CreatedDate,
		transaction.Version,
	).Scan(&transactionId, &createdAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	return returningUser, prevBalance, transactionId, createdAt, nil

}

func (b *BankRepository) UpdatePayment(ctx context.Context, user entity.User, remarks string) (entity.User, int, uuid.UUID, time.Time, error) {

	returningUser := entity.User{}
	tx, err := b.db.Begin(ctx)
	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	defer tx.Rollback(ctx)

	query := `update "user" set balance = balance - $1, version = version+1, updated_at = $2 where phone_number = $3 and version = $4 RETURNING id, balance, updated_at `

	selectUser := `select phone_number, balance, version from "user" where phone_number = $1`

	transactionQuery := `
		INSERT INTO transaction (id, amount, balance_before, balance_after, transaction_type, user_id, created_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at
	`

	var PhoneNumber string
	var prevBalance int
	var Version int

	err = tx.QueryRow(ctx, selectUser, user.PhoneNumber).Scan(&PhoneNumber, &prevBalance, &Version)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = pgsql.ErrUserNotFound
		}
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	if prevBalance < user.Balance {
		err = pgsql.ErrBalanceNotEnough
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	err = tx.QueryRow(ctx, query, user.Balance, time.Now(), PhoneNumber, Version).Scan(&returningUser.ID, &returningUser.Balance, &returningUser.UpdatedAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	id, err := pkg.GenerateId()
	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	// Sample transaction data
	transaction := entity.Transaction{
		ID:              id,
		Amount:          user.Balance,
		BalanceBefore:   prevBalance,
		Remarks:         remarks,
		BalanceAfter:    returningUser.Balance,
		TransactionType: "DEBIT",
		UserID:          returningUser.ID,
		CreatedDate:     time.Now(),
		Version:         1,
	}

	var transactionId uuid.UUID
	var createdAt time.Time

	err = tx.QueryRow(context.Background(), transactionQuery,
		transaction.ID,
		transaction.Amount,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.TransactionType,
		transaction.UserID,
		transaction.CreatedDate,
		transaction.Version,
	).Scan(&transactionId, &createdAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	return returningUser, prevBalance, transactionId, createdAt, nil
}

func (b *BankRepository) TransferTX(ctx context.Context, user entity.User, targetUser uuid.UUID, remarks string) (entity.User, int, uuid.UUID, time.Time, error) {
	returningUser := entity.User{}
	tx, err := b.db.Begin(ctx)
	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	defer tx.Rollback(ctx)

	updateOriginBalance := `update "user" set balance = balance - $1, version = version+1, updated_at = $2 where phone_number = $3 and version = $4 RETURNING id, balance, updated_at `

	selectUserOrigin := `select phone_number, balance, version from "user" where phone_number = $1`

	transactionQuery := `
		INSERT INTO transaction (id, amount, balance_before, balance_after, transaction_type, user_id, created_at, user_id_destination, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at
	`

	var PhoneNumberOrigin string
	var prevBalanceOrigin int
	var VersionOrigin int

	err = tx.QueryRow(ctx, selectUserOrigin, user.PhoneNumber).Scan(&PhoneNumberOrigin, &prevBalanceOrigin, &VersionOrigin)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = pgsql.ErrUserNotFound
		}
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	if prevBalanceOrigin < user.Balance {
		err = pgsql.ErrBalanceNotEnough
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	err = tx.QueryRow(ctx, updateOriginBalance, user.Balance, time.Now(), PhoneNumberOrigin, VersionOrigin).Scan(&returningUser.ID, &returningUser.Balance, &returningUser.UpdatedAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	// update destination user
	returningDestUser := entity.User{}

	updateQueryDestination := `update "user" set balance = balance + $1, version = version+1, updated_at = $2 where phone_number = $3 and version = $4 RETURNING id, balance, updated_at `

	selectUserDestination := `select phone_number, balance, version from "user" where id = $1`

	var PhoneNumberDestination string
	var prevBalanceDestination int
	var VersionDestination int

	err = tx.QueryRow(ctx, selectUserDestination, targetUser).Scan(&PhoneNumberDestination, &prevBalanceDestination, &VersionDestination)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = pgsql.ErrUserNotFound
		}
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	err = tx.QueryRow(ctx, updateQueryDestination, user.Balance, time.Now(), PhoneNumberDestination, VersionDestination).Scan(&returningDestUser.ID, &returningDestUser.Balance, &returningDestUser.UpdatedAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	id, err := pkg.GenerateId()
	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	// insert transaction
	// Sample transaction data
	transaction := entity.Transaction{
		ID:              id,
		Amount:          user.Balance,
		BalanceBefore:   prevBalanceOrigin,
		Remarks:         remarks,
		BalanceAfter:    returningUser.Balance,
		TransactionType: "DEBIT",
		TargetUserID:    returningDestUser.ID,
		UserID:          returningUser.ID,
		CreatedDate:     time.Now(),
		Version:         1,
	}

	var transactionId uuid.UUID
	var createdAt time.Time

	err = tx.QueryRow(context.Background(), transactionQuery,
		transaction.ID,
		transaction.Amount,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.TransactionType,
		transaction.UserID,
		transaction.CreatedDate,
		transaction.TargetUserID,
		transaction.Version,
	).Scan(&transactionId, &createdAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	return returningUser, prevBalanceOrigin, transactionId, createdAt, nil

}
