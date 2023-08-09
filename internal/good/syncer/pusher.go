package syncer

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"goods-service/internal/good/domain"
)

type LogWriter interface {
	WriteLog(ctx context.Context, log domain.Log) (err error)
}

type LogPusher struct {
	pool   *pgxpool.Pool
	writer LogWriter
}

func NewLogPusher() {}

func (p *LogPusher) 
