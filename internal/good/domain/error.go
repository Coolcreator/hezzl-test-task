package domain

import "errors"

var (
	ErrGoodNotFound = errors.New("good not found")
	ErrBadRequest   = errors.New("bad request")
)
