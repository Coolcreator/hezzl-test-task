package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"goods-service/internal/good/domain"
)

type LogWriter struct {
	conn    *nats.Conn
	subject string
}

func NewLogWriter(subject string) *LogWriter {
	return &LogWriter{
		subject: subject,
	}
}

func (w *LogWriter) SendLog(ctx context.Context, log domain.Log) (err error) {
	jsonData, err := json.Marshal(toNLog(log))
	if err != nil {
		err = fmt.Errorf("json marshal: %w", err)
		return err
	}

	err = w.conn.Publish(w.subject, jsonData)
	if err != nil {
		err = fmt.Errorf("publish log: %w", err)
		return
	}
	return
}

func toNLog(log domain.Log) (nlog nlog) {
	nlog.ID = log.ID
	nlog.ProjectID = log.ProjectID
	nlog.Name = log.Name
	nlog.Description = log.Description
	nlog.Priority = log.Priority
	nlog.Removed = log.Removed
	nlog.EventTime = log.EventTime
	return
}
