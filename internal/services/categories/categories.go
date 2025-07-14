package categories

import (
	"context"
	"fmt"
	"tradeservice/internal/models"
	"tradeservice/internal/storage"
)

type StorageCategories struct {
	storage storage.StoreRepository
}

func New(storage storage.StoreRepository) *StorageCategories {
	return &StorageCategories{
		storage: storage,
	}
}

func (c StorageCategories) Add(ctx context.Context, name string, productId string) (id string, err error) {

	id, err = c.storage.Add(ctx, name, productId)

	if err != nil {
		return id, fmt.Errorf("failed to add category %w", err)
	}

	return id, nil
}

func (c StorageCategories) Set(ctx context.Context, id string, name string) error {

	err := c.storage.Set(ctx, id, name)

	if err != nil {
		return fmt.Errorf("failed to set category %w", err)
	}

	return nil
}

func (c StorageCategories) Get(ctx context.Context) ([]models.Category, error) {
	category, err := c.storage.Get(ctx)

	if err != nil {
		return []models.Category{}, fmt.Errorf("failed to get categories %w", err)
	}

	return category, nil
}

func (c StorageCategories) Delete(ctx context.Context, id string) error {

	err := c.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete category %w", err)
	}

	return nil
}
