package repository

import (
	"context"
	"ecommerce/internal/order/entity"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetAll(ctx context.Context) ([]*entity.Order, error)
	GetByID(ctx context.Context, id int) (*entity.Order, error)
	GetUserOrders(ctx context.Context, username string) ([]*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id int) error
	GetInvoice(ctx context.Context, orderID int) ([]*entity.InvoiceData, error)
}
