package clickhouse

import (
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func NewConnection(addr string) (conn driver.Conn, err error) {
	options := &clickhouse.Options{Addr: []string{addr}}
	conn, err = clickhouse.Open(options)
	if err != nil {
		err = fmt.Errorf("clickhouse open: %v", err)
		return
	}
	return
}
