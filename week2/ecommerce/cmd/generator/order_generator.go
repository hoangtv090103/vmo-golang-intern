package main

import (
	"ecommerce/config"
	usercase "ecommerce/internal/faker_generator/usecase"
	orderDB "ecommerce/internal/order/infra/db"
	productDB "ecommerce/internal/product/infra/db"
	userDB "ecommerce/internal/user/infra/db"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize Postgres
	pg := config.NewPG()
	defer func(pg *config.PG) {
		err := pg.Close()
		if err != nil {
			return
		}
	}(pg)

	userRepo := userDB.NewUserRepoPG(pg, nil)
	orderRepo := orderDB.NewOrderRepoPG(pg)
	productRepo := productDB.NewProductRepoPG(pg)

	// Create a new generator
	generator := usercase.NewGeneratorUseCase(
		userRepo,
		orderRepo,
		productRepo,
	)

	// Generate data
	//generator.GenerateUser(userRepo)
	//generator.GenerateProduct(productRepo)
	generator.GenerateOrder(orderRepo)

}
