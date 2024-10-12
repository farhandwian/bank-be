package config

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-playground/validator/v10"
)

type UserConfig struct {
	PGx *pgxpool.Pool
	// Producer *sarama.SyncProducer
	Fiber    *fiber.App
	Validate *validator.Validate
}
