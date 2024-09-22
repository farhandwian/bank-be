package cmd

import (
	"bank-worker/feature/bank"
	"bank-worker/feature/shared"
	"bank-worker/pkg"
	"context"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func runTransferConsumer(ctx context.Context) {
	cfg := shared.LoadConfig("config/app.yml")
	kafkaCfg := pkg.NewKafkaConsumerConfig()

	consumer, err := sarama.NewConsumerGroup([]string{"172.17.0.1:9092"}, bank.CreateNewTransferTopic, kafkaCfg)
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

	bank.SetDBPool(pool)

	go func() {
		for err = range consumer.Errors() {
			log.Printf("consumer error, topic %s, error %s", bank.CreateNewTransferTopic, err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-newCtx.Done():
				log.Println("consumer stopped")
				return
			default:
				err = consumer.Consume(newCtx, []string{bank.CreateNewTransferTopic},
					pkg.NewKafkaConsumer(&bank.NewTransferEventHandler{}, 1000),
				)
				if err != nil {
					log.Printf("consume message error, topic %s, error %s", bank.CreateNewTransferTopic, err.Error())
					return
				}
			}
		}
	}()

	log.Printf("consumer up and running, topic %s, group: %s", bank.CreateNewTransferTopic, bank.CreateNewTransferTopicGroupConsumer)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	cancel()
	log.Println("cancelled message without marking offsets")
}
