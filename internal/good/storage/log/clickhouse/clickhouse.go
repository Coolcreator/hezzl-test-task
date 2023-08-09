package clickhouse

import (
	"context"

	"goods-service/internal/good/domain"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type LogStorage struct {
	conn driver.Conn
}

func NewLogStorage(conn driver.Conn) *LogStorage {
	return &LogStorage{}
}

func (s *LogStorage) WriteLogs(ctx context.Context, logs []domain.Log) (err error) {
	const query = `INSERT INTO logs`

	batch, err := s.conn.PrepareBatch(ctx, query)
	if err != nil {
		return
	}

	for _, log := range logs {
		err := batch.Append(log.ID, log.ProjectID, log.Name, log.Description, log.Priority, log.Removed, log.EventTime)
		if err != nil {
			return err
		}
	}

	batch.Send()

	return
}
