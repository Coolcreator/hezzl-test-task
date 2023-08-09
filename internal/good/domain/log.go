package domain

import "time"

type Log struct {
	ID          int64
	ProjectID   int64
	Name        string
	Description string
	Priority    int32
	Removed     bool
	EventTime   time.Time
}
