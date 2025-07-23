package storage

import (
	"context"
	"tradeservice/internal/models"
)

type CategoryRepository interface {
	AddCategory(ctx context.Context, name string, productId string) (id string, err error)
	GetCategory(ctx context.Context) ([]models.CategoryDto, error)
	SetCategory(ctx context.Context, id string, name string) error
	DeleteCategory(ctx context.Context, id string) error
}

type ProductRepository interface {
	AddProduct(ctx context.Context, name string) (id string, err error)
	GetProduct(ctx context.Context) ([]models.ProductDto, error)
	SetProduct(ctx context.Context, id string, name string) error
	DeleteProduct(ctx context.Context, id string) error
}
