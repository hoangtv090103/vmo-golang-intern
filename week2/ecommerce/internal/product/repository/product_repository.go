package repository

import "ecommerce/internal/product/domain"

type ProductRepository interface {
	Create(product domain.Product) error
	GetAll() ([]domain.Product, error)
	GetByID(id int) (domain.Product, error)
	GetByName(name string) ([]domain.Product, error)
	Update(product domain.Product) error
	Delete(id int) error
}
