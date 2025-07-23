package models

type CategoryDto struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ProductId string `json:"productId"`
}

type ProductDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
