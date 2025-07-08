package storage

import (
	"context"
	"tradeservice/internal/models"
)

type Repository interface {
	Add(ctx context.Context, name string, productId string) (id string, err error)
	Get(ctx context.Context) ([]models.Category, error)
	Set(ctx context.Context, id string, name string) error
	Delete(ctx context.Context, id string) error
}
