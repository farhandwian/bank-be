package repository

import (
	"context"

	"bank-backend/module/user/entity"
	"bank-backend/utils/pgsql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) InsertUser(ctx context.Context, user entity.User) (entity.User, error) {

	returningUser := entity.User{}

	query := `
        INSERT INTO "user" (id, phone_number, pin, first_name, last_name, address, created_at, updated_at, version, balance)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, first_name, last_name, phone_number, address,  created_at
    `
	err := u.db.QueryRow(ctx, query, user.ID, user.PhoneNumber, user.Pin, user.FirstName, user.LastName, user.Address, user.CreatedAt, user.UpdatedAt, user.Version, user.Balance).Scan(&returningUser.ID, &returningUser.FirstName, &returningUser.LastName, &returningUser.PhoneNumber, &returningUser.Address, &returningUser.CreatedAt)

	if err != nil {
		return returningUser, err
	}

	return returningUser, nil
}

func (u *UserRepository) CheckPhoneNumberExists(ctx context.Context, phoneNumber string) (bool, string, string, error) {
	query := `SELECT phone_number,pin FROM "user" WHERE phone_number = $1`
	var PhoneNumber string
	var Pin string

	err := u.db.QueryRow(ctx, query, phoneNumber).Scan(&PhoneNumber, &Pin)

	if err != nil {
		// Handle case where no rows are found
		if err == pgx.ErrNoRows {
			return false, "", "", nil
		}
		// Return any other error that occurs
		return false, "", "", err
	}
	// If a row is found, return true and the username
	return true, PhoneNumber, Pin, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {

	returningUser := entity.User{}
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return returningUser, err
	}
	defer tx.Rollback(ctx)

	query := `update "user" set first_name = $1 , last_name =$2, address= $3, version = version+1, updated_at = $4 where phone_number = $5 and version = $6 RETURNING id, first_name,last_name, address, updated_at `

	selectUser := `select phone_number, version from "user" where phone_number = $1`

	var PhoneNumber string
	var Version int

	err = tx.QueryRow(ctx, selectUser, user.PhoneNumber).Scan(&PhoneNumber, &Version)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = pgsql.ErrUserNotFound
		}
		return returningUser, err
	}

	err = tx.QueryRow(ctx, query, user.FirstName, user.LastName, user.Address, user.UpdatedAt, user.PhoneNumber, Version).Scan(&returningUser.ID, &returningUser.FirstName, &returningUser.LastName, &returningUser.Address, &returningUser.UpdatedAt)

	if err != nil {
		return returningUser, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return returningUser, err
	}

	return returningUser, nil

}
