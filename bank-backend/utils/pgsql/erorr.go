package pgsql

import "errors"

var (
	ErrUserNotFound     = errors.New("user: not found")
	ErrBalanceNotEnough = errors.New("bank: balance not enough")
)
