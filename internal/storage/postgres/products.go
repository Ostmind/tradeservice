package postgres

import (
	"context"
	"fmt"
	"tradeservice/internal/models"
)

type Products struct {
	db *Storage
}

func NewProducts(db *Storage) (*Categories, error) {
	return &Categories{
		db: db,
	}, nil
}

func (c *Products) Get(ctx context.Context) (product []models.Product, err error) {

	sqlStatement := `SELECT * FROM public.products`

	rows, err := c.db.DB.Query(ctx, sqlStatement)
	defer rows.Close()

	if err != nil {
		return product, fmt.Errorf("failed to query DB %w", err)
	}

	for rows.Next() {

		prod := models.Product{}

		err = rows.Scan(&prod.Id, &prod.Name, &prod.Created, &prod.Updated)

		if err != nil {
			return nil, fmt.Errorf("failed to parse DB %w", err)
		}
		product = append(product, prod)

	}

	return product, nil

}

func (c *Products) Add(ctx context.Context, name string, productId string) (id string, err error) {
	sqlStatement := `INSERT INTO public.products
					(name,product_id,created_at,updated_at) 
					values ($1,$2,now(),now());`

	_, err = c.db.DB.Exec(ctx, sqlStatement, name, productId)
	if err != nil {
		return "", fmt.Errorf("error adding to DB %w", err)
	}

	sqlStatement = `SELECT id FROM public.products where name = $1`

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

func (c *Products) Delete(ctx context.Context, id string) error {
	sqlStatement := `DELETE FROM public.products WHERE id = $1;`

	_, err := c.db.DB.Exec(ctx, sqlStatement, id)
	if err != nil {
		return fmt.Errorf("error deleting from DB %w", err)
	}

	return nil
}

func (c *Products) Set(ctx context.Context, id string, name string) error {
	sqlStatement := `UPDATE public.products SET name = $1 WHERE id = $2;`

	_, err := c.db.DB.Exec(ctx, sqlStatement, name, id)
	if err != nil {
		return fmt.Errorf("error updating DB %w", err)
	}

	return nil
}
