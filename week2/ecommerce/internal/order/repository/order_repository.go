package repository

import "ecommerce/internal/order/domain"

type OrderRepository interface {
	Create(order domain.Order) error
	GetAll() ([]domain.Order, error)
	GetByID(id int) (domain.Order, error)
	Update(order domain.Order) error
	Delete(id int) error
}
