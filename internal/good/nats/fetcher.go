package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type Fetcher struct {
	subscription *nats.Subscription
	batchSize    int32
}

func (f *Fetcher) FetchMessages(ctx context.Context) (err error) {
	messages, err := f.subscription.Fetch(int(f.batchSize))
	if err != nil {
		err = fmt.Errorf("fetch messages: %w", err)
		return
	}

	return
}

func NewFetcher(subscription *nats.Subscription) (fetcher *Fetcher) {
	fetcher = &Fetcher{
		subscription: subscription,
	}
	return
}
