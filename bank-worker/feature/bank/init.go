package bank

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db *pgxpool.Pool
)

func SetDBPool(dbPool *pgxpool.Pool) {
	if dbPool == nil {
		panic("cannot assign nil db pool")
	}

	db = dbPool
}
