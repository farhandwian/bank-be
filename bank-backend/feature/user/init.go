package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db       *pgxpool.Pool
	validate *validator.Validate
)

func SetDBPool(dbPool *pgxpool.Pool) {
	if dbPool == nil {
		panic("cannot assign nil db pool")
	}

	db = dbPool
}

func SetValidator(validator *validator.Validate) {
	if validator == nil {
		panic("cannot assign nil db pool")
	}
	validate = validator
}
