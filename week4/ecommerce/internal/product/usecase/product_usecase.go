package usecase

import (
	"context"
	"ecommerce/internal/product/entity"
	"ecommerce/internal/product/repository"
	"errors"
)

type ProductUsecase struct {
	productRepo repository.IProductRepository
}

func NewProductUsecase(productRepo repository.IProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

func (pu *ProductUsecase) CreateProduct(ctx context.Context, product *entity.Product) error {
    if product.Price < 0 {
		return errors.New("invalid price")
	}
	if product.Stock < 0 {
		return errors.New("invalid stock")
	}
	return pu.productRepo.Create(ctx, product)
}

func (pu *ProductUsecase) GetAllProducts(ctx context.Context) ([]*entity.Product, error) {
	return pu.productRepo.GetAll(ctx)
}

func (pu *ProductUsecase) GetByProductID(ctx context.Context, id int) (*entity.Product, error) {
	return pu.productRepo.GetByID(ctx, id)
}

func (pu *ProductUsecase) GetByProductName(ctx context.Context, name string) ([]*entity.Product, error) {
	return pu.productRepo.GetByName(ctx, name)
}

func (pu *ProductUsecase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	return pu.productRepo.Update(ctx, product)
}

func (pu *ProductUsecase) DeleteProduct(ctx context.Context, id int) error {
	return pu.productRepo.Delete(ctx, id)
}
