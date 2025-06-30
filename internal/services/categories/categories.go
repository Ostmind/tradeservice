package categories

import (
	"context"
	"fmt"
	"tradeservice/internal/models"
)

type Storager interface {
	Add(ctx context.Context, name string, productId string) (id string, err error)
	Get(ctx context.Context) (category []models.Category, err error)
	Set(ctx context.Context, id string, name string) error
	Delete(ctx context.Context, id string) error
}

type Categories struct {
	storage Storager
}

func New(storage Storager) *Categories {
	return &Categories{
		storage: storage,
	}
}

func (c Categories) Add(ctx context.Context, name string, productId string) (id string, err error) {

	id, err = c.storage.Add(ctx, name, productId)

	if err != nil {
		return id, fmt.Errorf("error query DB %w", err)
	}

	return id, nil
}

func (c Categories) Set(ctx context.Context, id string, name string) error {

	err := c.storage.Set(ctx, id, name)

	if err != nil {
		return fmt.Errorf("error updating DB %w", err)
	}

	return nil
}

func (c Categories) Get(ctx context.Context) ([]models.Category, error) {
	category, err := c.storage.Get(ctx)

	if err != nil {
		return []models.Category{}, fmt.Errorf("error quering DB %w", err)
	}

	return category, nil
}

func (c Categories) Delete(ctx context.Context, id string) error {

	err := c.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting DB %w", err)
	}

	return nil
}
