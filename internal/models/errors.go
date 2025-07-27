package models

import (
	"errors"
)

var (
	ErrUnique               = errors.New("already exists")
	ErrNotFound             = errors.New("not found")
	ErrDB                   = errors.New("db error")
	ErrDBConnectionCreation = errors.New("db connection creation error")
)
