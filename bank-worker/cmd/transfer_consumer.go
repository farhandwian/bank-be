package cmd

import (
	"book-service/feature/books"
	"book-service/feature/shared"
	"book-service/pkg"
	"context"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func runTransferConsumer(ctx context.Context) {
	cfg := shared.LoadConfig("config/created_user_consumer.yaml")
	kafkaCfg := pkg.NewKafkaConsumerConfig()

	consumer, err := sarama.NewConsumerGroup([]string{"172.17.0.1:9092"}, books.CreateNewUserTopicGroupConsumer, kafkaCfg)
	if err != nil {
		log.Fatalln("unable to create consumer group", err)
	}

	defer consumer.Close()

	dbCfg, err := pgxpool.ParseConfig(cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	// Set needed dependencies
	newCtx, cancel := context.WithCancel(ctx)

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln("unable to create database connection pool", err)
	}
	defer pool.Close()

	books.SetDBPool(pool)

	go func() {
		for err = range consumer.Errors() {
			log.Printf("consumer error, topic %s, error %s", books.CreateNewUserTopic, err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-newCtx.Done():
				log.Println("consumer stopped")
				return
			default:
				err = consumer.Consume(newCtx, []string{books.CreateNewUserTopic},
					pkg.NewKafkaConsumer(&books.NewUserEventHandler{}, 1000),
				)
				if err != nil {
					log.Printf("consume message error, topic %s, error %s", books.CreateNewUserTopic, err.Error())
					return
				}
			}
		}
	}()

	log.Printf("consumer up and running, topic %s, group: %s", books.CreateNewUserTopic, books.CreateNewUserTopicGroupConsumer)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	cancel()
	log.Println("cancelled message without marking offsets")
}
