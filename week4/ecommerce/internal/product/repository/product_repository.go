package repository

import (
	"context"
	"ecommerce/internal/product/entity"
)

type IProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetAll(ctx context.Context) ([]*entity.Product, error)
	GetByID(ctx context.Context, id int) (*entity.Product, error)
	GetByName(ctx context.Context, name string) ([]*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id int) error
}
