package config

import (
	bankcfg "bank-backend/module/bank/config"
	bank "bank-backend/module/bank/transport"
	usercfg "bank-backend/module/user/config"
	user "bank-backend/module/user/transport"
	"bank-backend/pkg"
	"bank-backend/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func StartHTTPServer(ctx context.Context) {
	// Load configuration
	cfg := LoadConfig("../../config/app.yml")

	userCfg := usercfg.UserConfig{}
	bankCfg := bankcfg.BankConfig{}
	// init db pool
	pool := InitializeDatabase(cfg.DBConfig, ctx)
	userCfg.PGx = pool
	bankCfg.PGx = pool

	defer pool.Close()

	validate := validator.New()
	validate.RegisterValidation("indonesianphone", utils.ValidateIndonesianPhoneNumber)
	userCfg.Validate = validate
	bankCfg.Validate = validate

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, pkg.NewKafkaProducerConfig())
	if err != nil {
		log.Fatalln("unable to create kafka producer", err)
	}

	defer producer.Close()
	bankCfg.Producer = &producer
	bankCfg.ProcessTranferTopic = cfg.ProcessTransferTopic

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

	userCfg.Fiber = app
	bankCfg.Fiber = app

	if app == nil {
		fmt.Println("testes1")
		log.Fatalln("Fiber app is not initialized")
	}

	fmt.Println("testes2")

	user.NewRest(userCfg)
	bank.NewRest(bankCfg)

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
