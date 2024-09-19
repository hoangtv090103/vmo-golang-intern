package domain

import (
	productDomain "ecommerce/internal/product/domain"
)

type OrderLine struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Qty       int     `json:"qty"`
	Total     float64 `json:"total"`

	Product productDomain.Product `json:"product"`
	//Order   Order                 `json:"order"`
}
