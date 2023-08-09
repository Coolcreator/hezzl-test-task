package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func NewConnection(url string) (conn *nats.Conn, err error) {
	conn, err = nats.Connect(url)
	if err != nil {
		err = fmt.Errorf("connect: %w", err)
		return
	}
	return
}
