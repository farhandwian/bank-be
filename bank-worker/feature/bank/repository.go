package bank

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

func transferTX(ctx context.Context, user User, targetUser uuid.UUID, remarks string, created time.Time, transferId string) (User, int, uuid.UUID, time.Time, error) {
	returningUser := User{}
	tx, err := db.Begin(ctx)
	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	defer tx.Rollback(ctx)

	updateOriginBalance := `update "user" set balance = balance - $1, version = version+1, updated_at = $2 where phone_number = $3 and version = $4 RETURNING id, balance, updated_at `

	selectUserOrigin := `select phone_number, balance, version from "user" where phone_number = $1`

	transactionQuery := `
		INSERT INTO transaction (id, amount, balance_before, balance_after, transaction_type, user_id, created_at, version, remarks)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at
	`

	var PhoneNumberOrigin string
	var prevBalanceOrigin int
	var VersionOrigin int

	err = tx.QueryRow(ctx, selectUserOrigin, user.PhoneNumber).Scan(&PhoneNumberOrigin, &prevBalanceOrigin, &VersionOrigin)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errUserNotFound
		}
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	err = tx.QueryRow(ctx, updateOriginBalance, user.Balance, time.Now(), PhoneNumberOrigin, VersionOrigin).Scan(&returningUser.ID, &returningUser.Balance, &returningUser.UpdatedAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	// update destination user
	returningDestUser := User{}

	updateQueryDestination := `update "user" set balance = balance + $1, version = version+1, updated_at = $2 where phone_number = $3 and version = $4 RETURNING id, balance, updated_at `

	selectUserDestination := `select phone_number, balance, version from "user" where id = $1`

	var PhoneNumberDestination string
	var prevBalanceDestination int
	var VersionDestination int

	err = tx.QueryRow(ctx, selectUserDestination, targetUser).Scan(&PhoneNumberDestination, &prevBalanceDestination, &VersionDestination)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = errUserNotFound
		}
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	err = tx.QueryRow(ctx, updateQueryDestination, user.Balance, time.Now(), PhoneNumberDestination, VersionDestination).Scan(&returningDestUser.ID, &returningDestUser.Balance, &returningDestUser.UpdatedAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	// insert transaction origin
	// Sample transaction data
	parse, err := uuid.Parse(transferId)
	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}
	transaction := Transaction{
		ID:              parse,
		Amount:          user.Balance,
		BalanceBefore:   prevBalanceOrigin,
		Remarks:         remarks,
		BalanceAfter:    returningUser.Balance,
		TransactionType: "DEBIT",
		UserID:          returningUser.ID,
		CreatedDate:     created,
		Version:         1,
	}

	transactionDestination := Transaction{
		ID:              parse,
		Amount:          user.Balance,
		BalanceBefore:   prevBalanceDestination,
		Remarks:         remarks,
		BalanceAfter:    returningDestUser.Balance,
		TransactionType: "CREDIT",
		UserID:          returningDestUser.ID,
		CreatedDate:     created,
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
		transaction.Remarks,
	).Scan(&transactionId, &createdAt)

	if err != nil {
		return returningUser, 0, uuid.UUID{}, time.Time{}, err
	}

	err = tx.QueryRow(context.Background(), transactionQuery,
		transactionDestination.ID,
		transactionDestination.Amount,
		transactionDestination.BalanceBefore,
		transactionDestination.BalanceAfter,
		transactionDestination.TransactionType,
		transactionDestination.UserID,
		transactionDestination.CreatedDate,
		transactionDestination.Version,
		transaction.Remarks,
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
