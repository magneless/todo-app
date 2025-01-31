package storage

import "errors"

var (
	ErrUserExists = errors.New("user exists")
)

const (
	UniqueViolationErrorCode = "23505"
)
