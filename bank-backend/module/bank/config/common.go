package config

import (
	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-playground/validator/v10"
)

type BankConfig struct {
	PGx                 *pgxpool.Pool
	Producer            *sarama.SyncProducer
	Fiber               *fiber.App
	Validate            *validator.Validate
	ProcessTranferTopic string
}
