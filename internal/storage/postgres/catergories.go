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

func (c *Categories) Get(ctx context.Context) (category []models.Category, err error) {
	var categoryDto []models.CategoryDto

	sqlStatement := `SELECT * FROM public.categories`

	rows, err := c.db.DB.Query(ctx, sqlStatement)
	defer rows.Close()

	if err != nil {
		return category, fmt.Errorf("failed to query DB %w", err)
	}

	for rows.Next() {

		cat := models.Category{}

		err = rows.Scan(&cat.Id, &cat.Name, &cat.Created, &cat.ProductId, &cat.Updated)

		if err != nil {
			return nil, fmt.Errorf("failed to parse DB %w", err)
		}
		category = append(category, cat)

		categoryDto = append(categoryDto, models.CategoryDto{Id: cat.Id, ProductId: cat.ProductId, Name: cat.Name})
	}

	return category, nil

}

func (c *Categories) Add(ctx context.Context, name string, productId string) (id string, err error) {
	sqlStatement := `INSERT INTO public.categories
					(id,product_id,created_at,updated_at) 
					values ($1,$2,now(),now());`

	result, err := c.db.DB.Exec(ctx, sqlStatement, name, productId)

	if err != nil {
		if result.Insert() == false {
			return "", models.ErrUnique
		}
		return "", fmt.Errorf("error adding to DB %w", err)
	}

	sqlStatement = `SELECT id FROM public.categories where name = $1`

	rows, err := c.db.DB.Query(ctx, sqlStatement)
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

func (c *Categories) Delete(ctx context.Context, id string) error {
	sqlStatement := `DELETE FROM public.categories WHERE id = $1;`

	_, err := c.db.DB.Exec(ctx, sqlStatement, id)
	if err != nil {
		return fmt.Errorf("error deleting from DB %w", err)
	}

	return nil
}

func (c *Categories) Set(ctx context.Context, id string, name string) error {
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
