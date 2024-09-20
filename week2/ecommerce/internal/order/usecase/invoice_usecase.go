package usecase

import (
	"ecommerce/internal/order/domain"
	orderDB "ecommerce/internal/order/infra/db"
	productDB "ecommerce/internal/product/infra/db"
	userDB "ecommerce/internal/user/infra/db"
)

type InvoiceService struct {
	OrderRepo   *orderDB.OrderRepoPG
	UserRepo    *userDB.UserRepoPG
	ProductRepo *productDB.ProductRepoPG
}

func NewInvoiceService(orderRepo *orderDB.OrderRepoPG, userRepo *userDB.UserRepoPG, productRepo *productDB.ProductRepoPG) *InvoiceService {
	return &InvoiceService{
		OrderRepo:   orderRepo,
		UserRepo:    userRepo,
		ProductRepo: productRepo,
	}
}

func (s *InvoiceService) GenerateInvoiceData(orderId int) (*domain.InvoiceData, error) {
	order, err := s.OrderRepo.GetByID(orderId)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepo.GetByID(order.UserID)
	if err != nil {
		return nil, err
	}

	var items []domain.InvoiceItem
	for _, line := range order.Lines {
		product, err := s.ProductRepo.GetByID(line.ProductID)
		if err != nil {
			return nil, err
		}
		items = append(items, domain.InvoiceItem{
			ProductName: product.Name,
			Quantity:    line.Qty,
			UnitPrice:   product.Price,
			TotalPrice:  line.Total,
		})
	}

	return &domain.InvoiceData{
		OrderID:      order.ID,
		OrderDate:    order.OrderDate.Format("2006-01-02"),
		CustomerName: user.Name,
		Items:        items,
		Total:        order.TotalPrice,
	}, nil
}
