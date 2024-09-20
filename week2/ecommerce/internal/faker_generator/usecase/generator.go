package usecase

import (
	orderDomain "ecommerce/internal/order/domain"
	db2 "ecommerce/internal/order/infra/db"
	productDomain "ecommerce/internal/product/domain"
	db3 "ecommerce/internal/product/infra/db"
	userDomain "ecommerce/internal/user/domain"
	"ecommerce/internal/user/infra/db"
	"github.com/go-faker/faker/v4"
	"math/rand"
	"time"
)

type IGenerator interface {
	GenerateOrder() []orderDomain.Order
	GenerateOrderLine() []orderDomain.OrderLine
	GenerateUser() []userDomain.User
	GenerateProduct() []productDomain.Product
}

type GeneratorUseCase struct {
	userRepo    db.UserRepoPG
	orderRepo   db2.OrderRepoPG
	productRepo db3.ProductRepoPG
}

func NewGeneratorUseCase(userRepo *db.UserRepoPG, orderRepo *db2.OrderRepoPG, productRepo *db3.ProductRepoPG) *GeneratorUseCase {
	return &GeneratorUseCase{
		userRepo:    *userRepo,
		orderRepo:   *orderRepo,
		productRepo: *productRepo,
	}
}

func (g *GeneratorUseCase) GenerateOrder(orderRepo *db2.OrderRepoPG) []orderDomain.Order {
	var orders []orderDomain.Order
	var order orderDomain.Order
	for ordNum := 1; ordNum <= 100; ordNum++ {
		orderDateStr := faker.Date()
		orderDate, err := time.Parse("2006-01-02", orderDateStr)

		if err != nil {
			panic(err)
		}

		// Create OrderLine
		lines := g.GenerateOrderLine(&g.orderRepo)

		// Create Order
		totalPrice := 0.0
		for _, line := range lines {
			totalPrice += line.Total
		}
		order = orderDomain.Order{
			UserID:     rand.Intn(10),
			OrderDate:  orderDate,
			TotalPrice: totalPrice,
			Lines:      lines,
		}

		orders = append(orders, order)
	}

	// Save to database
	for _, order := range orders {
		g.orderRepo.Create(order)
	}

	return orders

}

func (g *GeneratorUseCase) GenerateOrderLine(orderRepo *db2.OrderRepoPG) []orderDomain.OrderLine {
	var orderLines []orderDomain.OrderLine

	for ordLineNum := 1; ordLineNum <= rand.Intn(10)+1; ordLineNum++ {
		productID := rand.Intn(1000)
		qty := rand.Intn(10)
		// Fetch product from database
		product, err := g.productRepo.GetByID(productID)
		if err != nil {
			panic(err)
		}
		orderLine := orderDomain.OrderLine{
			Qty:       qty,
			ProductID: productID,
			Product:   product,
			Total:     product.Price * float64(qty),
		}

		orderLines = append(orderLines, orderLine)
	}

	return orderLines
}

func (g *GeneratorUseCase) GenerateProduct(productRepo *db3.ProductRepoPG) []productDomain.Product {
	var products []productDomain.Product
	var product productDomain.Product

	for prodNum := 1; prodNum <= 1000; prodNum++ {
		product = productDomain.Product{
			Name:        faker.Word(),
			Description: faker.Sentence(),
			Price:       rand.Float64(),
			Stock:       rand.Intn(100),
		}

		products = append(products, product)

		// Save to database
		productRepo.Create(product)
	}

	return products
}

func (g *GeneratorUseCase) GenerateUser(userRepo *db.UserRepoPG) []userDomain.User {
	var users []userDomain.User
	var user userDomain.User
	for useNum := 1; useNum <= 10; useNum++ {
		user = userDomain.User{
			Name:     faker.Name(),
			Username: faker.Username(),
			Email:    faker.Email(),
			Balance:  rand.Float64(),
		}

		users = append(users, user)
	}

	// Save to database
	for _, user := range users {
		userRepo.Create(user)
	}

	return users
}
