package product

import (
	"context"
	"fmt"
	"tradeservice/internal/models"
	"tradeservice/internal/storage"
)

type StorageProducts struct {
	storage storage.StoreRepository
}

func New(storage storage.StoreRepository) *StorageProducts {
	return &StorageProducts{
		storage: storage,
	}
}

func (c StorageProducts) Add(ctx context.Context, name string, productId string) (id string, err error) {

	id, err = c.storage.Add(ctx, name, productId)

	if err != nil {
		return id, fmt.Errorf("failed to add product %w", err)
	}

	return id, nil
}

func (c StorageProducts) Set(ctx context.Context, id string, name string) error {

	err := c.storage.Set(ctx, id, name)

	if err != nil {
		return fmt.Errorf("failed to set product %w", err)
	}

	return nil
}

func (c StorageProducts) Get(ctx context.Context) ([]models.Category, error) {
	category, err := c.storage.Get(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get product %w", err)
	}

	return category, nil
}

func (c StorageProducts) Delete(ctx context.Context, id string) error {

	err := c.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product %w", err)
	}

	return nil
}
