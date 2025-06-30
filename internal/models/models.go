package models

import "time"

type Category struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	ProductId string    `db:"productId"`
	Created   time.Time `db:"createdAt"`
	Updated   time.Time `db:"updatedAt"`
}
