package entity

import (
	"ecommerce/internal/product/entity"
)

type OrderLine struct {
	ID        int     `json:"id,omitempty"`
	OrderID   int     `json:"order_id,omitempty"`
	ProductID int     `json:"product_id,omitempty"`
	Qty       int     `json:"qty,omitempty"`
	Total     float64 `json:"total,omitempty"`

	Product entity.Product `json:"product,omitempty"`
	Order   Order          `json:"order,omitempty"`
}
