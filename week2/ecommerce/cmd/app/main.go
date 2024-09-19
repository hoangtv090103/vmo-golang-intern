package main

import (
	"context"
	"ecommerce/config"
	userHandler "ecommerce/internal/user/handler"
	userDB "ecommerce/internal/user/infra/db"
	userUsecase "ecommerce/internal/user/usecase"

	productHandler "ecommerce/internal/product/handler"
	productDB "ecommerce/internal/product/infra/db"
	productUsecase "ecommerce/internal/product/usecase"

	orderHandler "ecommerce/internal/order/handler"
	orderDB "ecommerce/internal/order/infra/db"
	orderUsecase "ecommerce/internal/order/usecase"

	"ecommerce/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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

	// Initialize Redis
	redis := config.NewRedis()
	defer redis.Close()

	// Initialize Context
	ctx := context.Background()

	// User
	userRepo := userDB.NewUserRepoPG(pg, redis, ctx)
	userUC := userUsecase.NewUserUseCase(userRepo)
	uh := userHandler.NewUserHandler(userUC)
	routes.UserRoute(uh)

	// Product
	productRe := productDB.NewProductRepoPG(pg)
	productUC := productUsecase.NewProductUseCase(productRe)
	pH := productHandler.NewProductHandler(productUC)
	routes.ProductRoute(pH)

	// Order
	orderRepo := orderDB.NewOrderRepoPg(pg)
	orderUC := orderUsecase.NewOrderUseCase(orderRepo)
	orderH := orderHandler.NewOrderHandler(orderUC)
	routes.OrderRoute(orderH)

	log.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
