package nats

import "time"

type nlog struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int32     `json:"priority"`
	Removed     bool      `json:"removed"`
	EventTime   time.Time `json:"event_at"`
}
