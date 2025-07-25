package postgres

import (
	"context"
	"fmt"
	"tradeservice/internal/models"
)

type Categories struct {
	db *Storage
}

func NewCategories(db *Storage) (*Categories, error) {
	return &Categories{
		db: db,
	}, nil
}

func (c *Categories) GetCategory(ctx context.Context) (categoryDto []models.CategoryDto, err error) {
	sqlStatement := `SELECT * FROM public.categories`

	rows, err := c.db.DB.Query(ctx, sqlStatement)
	if err != nil {
		return categoryDto, fmt.Errorf("failed to query DB %w", err)
	}

	defer rows.Close()

	if err != nil {
		return categoryDto, fmt.Errorf("failed to query DB %w", err)
	}

	for rows.Next() {
		cat := models.Category{}

		err = rows.Scan(&cat.ID, &cat.Name, &cat.Created, &cat.ProductID, &cat.Updated)

		if err != nil {
			return nil, fmt.Errorf("failed to parse DB %w", err)
		}

		categoryDto = append(categoryDto, models.CategoryDto{ID: cat.ID, ProductID: cat.ProductID, Name: cat.Name})
	}

	return categoryDto, nil
}

func (c *Categories) AddCategory(ctx context.Context, name string, productID string) (id string, err error) {
	sqlStatement := `INSERT INTO public.categories
					(name,product_id,created_at,updated_at) 
					values ($1,$2,now(),now());`

	result, err := c.db.DB.Exec(ctx, sqlStatement, name, productID)

	if err != nil {
		if !result.Insert() {
			return "", models.ErrUnique
		}

		return "", fmt.Errorf("error adding to DB %w", err)
	}

	sqlStatement = `SELECT id FROM public.categories where name = $1`

	rows, err := c.db.DB.Query(ctx, sqlStatement)
	if err != nil {
		return "", fmt.Errorf("failed to query DB %w", err)
	}

	defer rows.Close()

	if err != nil {
		return "", fmt.Errorf("failed to query DB %w", err)
	}

	err = rows.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to parse DB %w", err)
	}

	return id, nil
}

func (c *Categories) DeleteCategory(ctx context.Context, id string) error {
	sqlStatement := `DELETE FROM public.categories WHERE id = $1;`

	result, err := c.db.DB.Exec(ctx, sqlStatement, id)
	if err != nil {
		return fmt.Errorf("error deleting from DB %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (c *Categories) SetCategory(ctx context.Context, id string, name string) error {
	sqlStatement := `UPDATE public.categories SET name = $1 WHERE id = $2;`

	result, err := c.db.DB.Exec(ctx, sqlStatement, name, id)
	if err != nil {
		return fmt.Errorf("error updating DB %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
