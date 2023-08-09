package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"goods-service/internal/good/domain"
)

type Decoder interface {
	DecodeLogs(ctx context.Context, data [][]byte)
}

type LogReader struct {
	subscription *nats.Subscription
	batchSize    int32
}

func NewLogReader(subscription *nats.Subscription, batchSize int32) *LogReader {
	return &LogReader{
		subscription: subscription,
		batchSize:    batchSize,
	}
}

func (r *LogReader) FetchLogs(ctx context.Context) (logs []domain.Log, ackBatch func() (err error), err error) {
	messages, err := r.subscription.Fetch(int(r.batchSize))
	if err != nil {
		err = fmt.Errorf("subscription fetch: %w", err)
		return
	}
	nlogs := make([]nlog, len(messages))
	for _, message := range messages {
		var nlog nlog
		err = json.Unmarshal(message.Data, nlog)
		if err != nil {
			err = fmt.Errorf("json unmarshal: %w", err)
			return
		}
		nlogs = append(nlogs, nlog)
	}
	logs = toLogs(nlogs)
	ackBatch = func() (err error) {
		return messages[len(messages)-1].Ack()
	}
	return
}

func toLogs(nlogs []nlog) (logs []domain.Log) {
	logs = make([]domain.Log, 0, len(nlogs))
	for _, nlog := range nlogs {
		logs = append(logs, domain.Log{
			ID:          nlog.ID,
			ProjectID:   nlog.ProjectID,
			Name:        nlog.Name,
			Description: nlog.Description,
			Priority:    nlog.Priority,
			EventTime:   nlog.EventTime,
		})
	}
	return
}
