package domain

import (
	userDomain "ecommerce/internal/user/domain"
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	OrderDate  time.Time `json:"order_date"`
	TotalPrice float64   `json:"total_price"`

	User  userDomain.User `json:"user"`
	Lines []OrderLine     `json:"lines"`
}
