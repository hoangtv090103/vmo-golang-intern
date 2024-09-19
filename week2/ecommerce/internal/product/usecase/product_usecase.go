package usecase

import (
	productDomain "ecommerce/internal/product/domain"
	"ecommerce/internal/product/repository"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

func (uc *ProductUseCase) AddProduct(product productDomain.Product) error {
	return uc.productRepo.Create(product)
}

func (uc *ProductUseCase) GetAllProducts() ([]productDomain.Product, error) {
	return uc.productRepo.GetAll()
}

func (uc *ProductUseCase) GetProductByID(id int) (productDomain.Product, error) {
	return uc.productRepo.GetByID(id)
}

func (uc *ProductUseCase) UpdateProduct(product productDomain.Product) error {
	return uc.productRepo.Update(product)
}

func (uc *ProductUseCase) DeleteProduct(id int) error {
	return uc.productRepo.Delete(id)
}


