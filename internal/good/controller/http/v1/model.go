package v1

import "time"

type createGoodRequest struct {
	Name string `json:"name"`
}

type updateGoodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type reprioritizeRequest struct {
	NewPriority int32 `json:"newPriority"`
}

type goodPriority struct {
	ID       int64 `json:"id"`
	Priority int32 `json:"priority"`
}

type reprioritizeResult struct {
	Priotities []goodPriority `json:"priorities"`
}

type removeGoodResponse struct {
	ID        int64 `json:"id"`
	ProjectID int64 `json:"projectId"`
	Removed   bool  `json:"removed"`
}

type meta struct {
	Total   int32 `json:"total"`
	Removed int32 `json:"removed"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

type goodResult struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int32     `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

type goodsListResult struct {
	Meta  meta         `json:"meta"`
	Goods []goodResult `json:"goods"`
}
