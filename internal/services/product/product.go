package product

import (
	"context"
	"fmt"
	"tradeservice/internal/models"
	"tradeservice/internal/storage"
)

type StorageProducts struct {
	storage storage.Repository
}

func New(storage storage.Repository) *StorageProducts {
	return &StorageProducts{
		storage: storage,
	}
}

func (c StorageProducts) Add(ctx context.Context, name string, productId string) (id string, err error) {

	id, err = c.storage.Add(ctx, name, productId)

	if err != nil {
		return id, fmt.Errorf("failed to add category %w", err)
	}

	return id, nil
}

func (c StorageProducts) Set(ctx context.Context, id string, name string) error {

	err := c.storage.Set(ctx, id, name)

	if err != nil {
		return fmt.Errorf("failed to set category %w", err)
	}

	return nil
}

func (c StorageProducts) Get(ctx context.Context) ([]models.Category, error) {
	category, err := c.storage.Get(ctx)

	if err != nil {
		return []models.Category{}, fmt.Errorf("failed to get categories %w", err)
	}

	return category, nil
}

func (c StorageProducts) Delete(ctx context.Context, id string) error {

	err := c.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete category %w", err)
	}

	return nil
}
