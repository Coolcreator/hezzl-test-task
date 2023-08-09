package redis

import "time"

type (
	meta struct {
		Total   int32 `json:"total"`
		Removed int32 `json:"removed"`
		Limit   int32 `json:"limit"`
		Offset  int32 `json:"offset"`
	}

	good struct {
		ID          int64     `json:"id"`
		ProjectID   int64     `json:"projectId"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Priority    int32     `json:"priority"`
		Removed     bool      `json:"removed"`
		CreatedAt   time.Time `json:"created_at"`
	}

	listOfGoods struct {
		Meta  meta   `json:"meta"`
		Goods []good `json:"goods"`
	}
)
