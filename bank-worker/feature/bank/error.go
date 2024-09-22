package bank

import "errors"

var (
	errUserNotFound     = errors.New("user: not found")
	errBalanceNotEnough = errors.New("bank: balance not enough")
)
