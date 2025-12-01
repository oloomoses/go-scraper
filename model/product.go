package model

type Product struct {
	ID    uint   `json:"product_id"`
	Title string `json:"title"`
	Price string `json:"price"`
	Url   string `json:"URL"`
}
