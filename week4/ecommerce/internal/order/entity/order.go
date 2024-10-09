package entity

import (
	userDomain "ecommerce/internal/user/entity"
	"time"
)

type Order struct {
	ID         int        `json:"id,omitempty"`
	UserID     int        `json:"user_id,omitempty"`
	OrderDate  *time.Time `json:"created_at,omitempty"`
	TotalPrice float64    `json:"total_price,omitempty"`

	User  userDomain.User `json:"user,omitempty"`
	Lines []OrderLine     `json:"lines,omitempty"`
}
