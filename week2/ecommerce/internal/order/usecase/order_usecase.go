package usecase

import (
	"ecommerce/internal/order/domain"
	"ecommerce/internal/order/infra/db"
)

type OrderUsecase interface {
	CreateOrder(order domain.Order) error
	GetAllOrders() ([]domain.Order, error)
	GetOrderById(id int) (domain.Order, error)
	UpdateOrder(order domain.Order) error
	DeleteOrder(id int) error
}

type orderUsecase struct {
	orderRepo db.OrderRepoPG
}

func NewOrderUseCase(orderRepo *db.OrderRepoPG) OrderUsecase {
	return &orderUsecase{
		orderRepo: *orderRepo,
	}
}

func (uc *orderUsecase) CreateOrder(order domain.Order) error {
	return uc.orderRepo.Create(order)
}

func (uc *orderUsecase) GetAllOrders() ([]domain.Order, error) {
	return uc.orderRepo.GetAll()
}

func (uc *orderUsecase) GetOrderById(id int) (domain.Order, error) {
	return uc.orderRepo.GetByID(id)
}

func (uc *orderUsecase) UpdateOrder(order domain.Order) error {
	return uc.orderRepo.Update(order)
}

func (uc *orderUsecase) DeleteOrder(id int) error {
	return uc.orderRepo.Delete(id)
}
