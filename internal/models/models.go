package models

import "time"

type Category struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	ProductID string    `db:"product_id"`
	Created   time.Time `db:"created_at"`
	Updated   time.Time `db:"updated_at"`
}

type Product struct {
	ID      string    `db:"id"`
	Name    string    `db:"name"`
	Created time.Time `db:"created_at"`
	Updated time.Time `db:"updated_at"`
}
