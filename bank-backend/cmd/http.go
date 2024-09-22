package cmd

import (
	"bank-backend/feature/bank"
	"bank-backend/feature/shared"
	"bank-backend/feature/user"
	"bank-backend/pkg"
	"context"
	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func runHTTPServer(ctx context.Context) {
	// Load configuration
	cfg := shared.LoadConfig("config/app.yml")

	dbCfg, err := pgxpool.ParseConfig(cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	// Set needed dependencies
	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln("unable to create database connection pool", err)
	}
	defer pool.Close()
	user.SetDBPool(pool)
	bank.SetDBPool(pool)

	validate := validator.New()
	validate.RegisterValidation("indonesianphone", shared.ValidateIndonesianPhoneNumber)
	user.SetValidator(validate)
	bank.SetValidator(validate)

	producer, err := sarama.NewSyncProducer([]string{"172.17.0.1:9092"}, pkg.NewKafkaProducerConfig())
	if err != nil {
		log.Fatalln("unable to create kafka producer", err)
	}

	defer producer.Close()
	bank.SetKafkaProducer(producer)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	})

	// Health check route

	app.Get("/health", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	user.FiberRoute(app)
	bank.FiberRoute(app)
	go func() {

		if err := app.Listen(cfg.Server.Addr()); err != nil {
			log.Fatalln("unable to start server", err)
		}
	}()

	log.Println("server started")

	// Wait for signal to shut down
	<-ctx.Done()

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Fatalln("unable to shutdown server", err)
	}

	log.Println("server shutdown")
}
