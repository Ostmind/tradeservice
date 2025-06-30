package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
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
	var cat models.Category

	sqlStatement := `SELECT * FROM public.categories`

	rows, err := c.db.DB.Query(ctx, sqlStatement) //creating Struct of Data
	defer rows.Close()

	if err != nil {
		return category, fmt.Errorf("error query DB %w", err)
	}

	for rows.Next() {

		err = rows.Scan(&cat.Id, &cat.Name, &cat.Created, &cat.ProductId, &cat.Updated)

		if err != nil {
			return nil, fmt.Errorf("error parsing DB %w", err)
		}
		category = append(category, cat)

	}

	return category, nil

}

func (c *Categories) Add(ctx context.Context, name string, productId string) (id string, err error) {
	sqlStatement := `INSERT INTO public.categories(id,name,productId,createdAt,updatedAt) values ($1,$2,$3,now(),now());`

	id = "1-" + strconv.Itoa(rand.Intn(1000000))

	_, err = c.db.DB.Exec(ctx, sqlStatement, id, name, productId)
	if err != nil {
		return "0", fmt.Errorf("error adding to DB %w", err)
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

	_, err := c.db.DB.Exec(ctx, sqlStatement, name, id)
	if err != nil {
		return fmt.Errorf("error updating DB %w", err)
	}

	return nil
}
