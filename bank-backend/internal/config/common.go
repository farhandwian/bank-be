package config

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Server               serverConfig `yaml:"server" json:"server"`
	DBConfig             pgConfig     `yaml:"db" json:"db"`
	Kafka                kafkaConfig  `yaml:"kafka" json:"kafka"`
	ProcessTransferTopic string       `yaml:"process_transfer_topic" json:"process_transfer_topic"`
}

func loadConfigFromReader(r io.Reader, c *config) error {
	return yaml.NewDecoder(r).Decode(c)
}

func loadConfigFromFile(fn string, c *config) error {
	_, err := os.Stat(fn)

	if err != nil {
		return err
	}

	f, err := os.Open(fn)

	if err != nil {
		return err
	}

	defer f.Close()

	return loadConfigFromReader(f, c)
}

func LoadConfig(fn string) config {
	cfg := config{}
	err := loadConfigFromFile(fn, &cfg)
	if err != nil {
		panic(err)
	}

	slog.Debug("config loaded", slog.Any("config", cfg))
	return cfg
}
