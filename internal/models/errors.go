package models

import (
	"errors"
)

var (
	ErrUnique   = errors.New("already exists")
	ErrNotFound = errors.New("not found")
)
