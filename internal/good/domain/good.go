package domain

import "time"

type Good struct {
	ID          int64
	ProjectID   int64
	Name        string
	Description string
	Priority    int32
	Removed     bool
	CreatedAt   time.Time
}

type CreateGood struct {
	ProjectID int64
	Name      string
}

type UpdateGood struct {
	ID          int64
	ProjectID   int64
	Name        string
	Description string
}

type DeleteGood struct {
	ID        int64
	ProjectID int64
}

type ReprioritizeGood struct {
	ID          int64
	ProjectID   int64
	NewPriority int32
}

type GoodPriority struct {
	ID       int64
	Priority int32
}

type ListGoods struct {
	Limit  int32
	Offset int32
}

type Meta struct {
	Total   int32
	Removed int32
	Limit   int32
	Offset  int32
}

type GoodsList struct {
	Goods []Good
	Meta  Meta
}
