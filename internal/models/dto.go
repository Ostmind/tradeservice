package models

type CategoryDto struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProductID string `json:"productId"`
}

type ProductDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
