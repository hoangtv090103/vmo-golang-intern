package repository

import "ecommerce/internal/order/domain"

type OrderLineRepository interface {
	Create(orderLine domain.OrderLine) error
	GetAll() ([]domain.OrderLine, error)
	GetByID(id int) (domain.OrderLine, error)
	Update(orderLine domain.OrderLine) error
	Delete(id int) error
}