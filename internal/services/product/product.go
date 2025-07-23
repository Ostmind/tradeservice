package product

import (
	"context"
	"fmt"
	"tradeservice/internal/models"
	"tradeservice/internal/storage"
)

type StorageProducts struct {
	storage storage.ProductRepository
}

func New(storage storage.ProductRepository) *StorageProducts {
	return &StorageProducts{
		storage: storage,
	}
}

func (c StorageProducts) AddProduct(ctx context.Context, name string) (id string, err error) {

	id, err = c.storage.AddProduct(ctx, name)

	if err != nil {
		return id, fmt.Errorf("failed to add product %w", err)
	}

	return id, nil
}

func (c StorageProducts) SetProduct(ctx context.Context, id string, name string) error {

	err := c.storage.SetProduct(ctx, id, name)

	if err != nil {
		return fmt.Errorf("failed to set product %w", err)
	}

	return nil
}

func (c StorageProducts) GetProduct(ctx context.Context) ([]models.ProductDto, error) {
	product, err := c.storage.GetProduct(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get product %w", err)
	}

	return product, nil
}

func (c StorageProducts) DeleteProduct(ctx context.Context, id string) error {

	err := c.storage.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product %w", err)
	}

	return nil
}
