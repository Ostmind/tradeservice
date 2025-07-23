package models

import "time"

type Category struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	ProductId string    `db:"product_id"`
	Created   time.Time `db:"created_at"`
	Updated   time.Time `db:"updated_at"`
}

type Product struct {
	Id      string    `db:"id"`
	Name    string    `db:"name"`
	Created time.Time `db:"created_at"`
	Updated time.Time `db:"updated_at"`
}
