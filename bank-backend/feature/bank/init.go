package bank

import (
	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db       *pgxpool.Pool
	validate *validator.Validate
	kp       sarama.SyncProducer
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

func SetKafkaProducer(producer sarama.SyncProducer) {
	if producer == nil {
		panic("cannot assign nil kafka producer")
	}

	kp = producer
}
