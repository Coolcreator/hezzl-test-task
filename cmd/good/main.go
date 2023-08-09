package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"golang.org/x/exp/slog"

	"goods-service/internal/good/cache/redis"
	v1 "goods-service/internal/good/controller/http/v1"
	"goods-service/internal/good/service"
	"goods-service/internal/good/storage/good/postgres"

	cc "goods-service/pkg/clickhouse"
	hs "goods-service/pkg/http"
	ls "goods-service/pkg/log/slog"
	nc "goods-service/pkg/nats"
	pc "goods-service/pkg/postgres"
	rc "goods-service/pkg/redis"
)

type Config struct {
	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"debug"`
	}
	Postgres struct {
		User     string `env:"POSTGRES_USER" env-required:"true"`
		Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
		Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"POSTGRES_PORT" env-default:"5432"`
		DB       string `env:"POSTGRES_DB" env-required:"true"`
		SSLMode  string `env:"POSTGRES_SSL_MODE" env-default:"false"`
	}
	Clickhouse struct {
		Address string `env:"CLICKHOUSE_ADDRESS" env-required:"true"`
	}
	NATS struct {
		URL string `env:"NATS_URL" env-required:"true"`
	}
	Redis struct {
		URL string `env:"REDIS_URL" env-required:"true"`
	}
	HTTP struct{}
}

func main() {
	var (
		cfg Config
		log *slog.Logger
		err error
	)
	flag.Parse()
	log = ls.NewLogger(cfg.Log.Level)
	log.Info("starting good service...")
	log.Info("reading config...")
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Error("failed to read env", ls.Error(err))
		os.Exit(1)
	}
	log.Info("initializing clients...")
	redisClient, err := rc.NewClient(cfg.Redis.URL)
	if err != nil {
		log.Error("failed to create redis client", ls.Error(err))
		os.Exit(1)
	}
	defer redisClient.Close()
	natsConn, err := nc.NewConnection(cfg.NATS.URL)
	if err != nil {
		log.Error("failed to establish nats connection", ls.Error(err))
		os.Exit(1)
	}
	defer natsConn.Drain()
	pool, err := pc.NewConnPool(&pc.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DB:       cfg.Postgres.DB,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		log.Error("failed to create postgres connections pool", ls.Error(err))
		return
	}
	defer pool.Close()
	clickhouseConn, err := cc.NewConnection(cfg.Clickhouse.Address)
	if err != nil {
		log.Error("failed to establish clickhouse connection", ls.Error(err))
		os.Exit(1)
	}
	defer clickhouseConn.Close()
	cache := redis.NewCache(redisClient)
	storage := postgres.NewGoodStorage(pool)
	service := service.NewGoodService(cache, storage)
	controller := v1.NewController(service)
	mux := chi.NewRouter()
	controller.Register(mux)
	server := hs.NewServer(mux)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server.Run()

	<-ctx.Done()

	err = server.Shutdown()
	if err != nil {
		log.Error("failed to shutdown server", ls.Error(err))
	}
}
