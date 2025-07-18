package categories

import (
	"context"
	"fmt"
	"tradeservice/internal/models"
	"tradeservice/internal/storage"
)

type StorageCategories struct {
	storage storage.CategoryRepository
}

func New(storage storage.CategoryRepository) *StorageCategories {
	return &StorageCategories{
		storage: storage,
	}
}

func (c StorageCategories) AddCategory(ctx context.Context, name string, productId string) (id string, err error) {

	id, err = c.storage.AddCategory(ctx, name, productId)

	if err != nil {
		return id, fmt.Errorf("failed to add category %w", err)
	}

	return id, nil
}

func (c StorageCategories) SetCategory(ctx context.Context, id string, name string) error {

	err := c.storage.SetCategory(ctx, id, name)

	if err != nil {
		return fmt.Errorf("failed to set category %w", err)
	}

	return nil
}

func (c StorageCategories) GetCategory(ctx context.Context) ([]models.CategoryDto, error) {
	category, err := c.storage.GetCategory(ctx)

	if err != nil {
		return []models.CategoryDto{}, fmt.Errorf("failed to get categories %w", err)
	}

	return category, nil
}

func (c StorageCategories) DeleteCategory(ctx context.Context, id string) error {

	err := c.storage.DeleteCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete category %w", err)
	}

	return nil
}
