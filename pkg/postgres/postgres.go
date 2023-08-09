package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConn           = 10
	minConns          = 5
	healthCheckPeriod = 3 * time.Minute
	maxConnIdleTime   = 1 * time.Minute
	maxConnLifetime   = 3 * time.Minute
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	SSLMode  string
}

func NewConnPool(config *Config) (pool *pgxpool.Pool, err error) {
	connString := buildConnString(config)
	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		err = fmt.Errorf("parse config")
		return
	}
	poolCfg.MaxConns = maxConn
	poolCfg.MinConns = minConns
	poolCfg.HealthCheckPeriod = healthCheckPeriod
	poolCfg.MaxConnIdleTime = maxConnIdleTime
	poolCfg.MaxConnLifetime = maxConnLifetime
	ctx := context.Background()
	pool, err = pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		err = fmt.Errorf("new with config")
		return
	}
	return
}

func buildConnString(config *Config) (connString string) {
	connString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DB,
		config.SSLMode,
	)
	return
}
