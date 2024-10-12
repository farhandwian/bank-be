package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     uint   `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	DBName   string `yaml:"db_name" json:"db_name"`
	SslMode  string `yaml:"ssl_mode" json:"ssl_mode"`
	MinConn  uint   `yaml:"min_conn" json:"min_conn"`
	MaxConn  uint   `yaml:"max_conn" json:"max_conn"`
}

func InitializeDatabase(envConfig pgConfig, ctx context.Context) *pgxpool.Pool {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=%s pool_min_conns=%d pool_max_conns=%d",
		envConfig.User, envConfig.Password, envConfig.Host, envConfig.Port, envConfig.DBName, envConfig.SslMode, envConfig.MinConn, envConfig.MaxConn,
	)
	dbCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln("unable to create database connection pool", err)
	}

	return pool

}
