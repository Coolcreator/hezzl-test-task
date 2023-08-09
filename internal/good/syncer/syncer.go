package syncer

import (
	"context"

	"goods-service/internal/good/domain"
)

type LogReader interface {
	ReadLogs(ctx context.Context) (logs []domain.Log, err error)
}

type LogStorage interface {
	WriteLogs(ctx context.Context, logs []domain.Log) (err error)
}

type LogSyncer struct {
	reader  LogReader
	storage LogStorage
}

func NewLogSyncer(reader LogReader, storage LogStorage) *LogSyncer {
	return &LogSyncer{
		reader:  reader,
		storage: storage,
	}
}

func (s *LogSyncer) SyncLogs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logs, err := s.reader.ReadLogs(ctx)
			if err != nil {
				continue
			}
			for {
				err = s.storage.WriteLogs(ctx, logs)
				if err != nil {
					continue
				}

			}
		}
	}
}
