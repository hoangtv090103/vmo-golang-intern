// main.go
package main

import (
	//"context"
	"ecommerce/config"
	authHandler "ecommerce/internal/auth/handler"
	authDB "ecommerce/internal/auth/infra/db"
	authUsecase "ecommerce/internal/auth/usecase"

	productHandler "ecommerce/internal/product/handler"
	productDB "ecommerce/internal/product/infra/db"
	productUsecase "ecommerce/internal/product/usecase"
	userHandler "ecommerce/internal/user/handler"
	userDB "ecommerce/internal/user/infra/db"
	userUsecase "ecommerce/internal/user/usecase"
	"ecommerce/middleware"

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
	//ctx := context.Background()

	// Add middleware
	mux := http.NewServeMux()

	// Create a new ServeMux for auth routes
	authMux := http.NewServeMux()
	// Auth
	authRepo := authDB.NewAuthRepoPG(pg)
	authUC := authUsecase.NewAuthUsecase(authRepo)
	authH := authHandler.NewAuthHandler(authUC)
	routes.AuthRoute(authH, authMux)

	// User
	userRepo := userDB.NewUserRepoPG(pg, redis)
	userUC := userUsecase.NewUserUseCase(userRepo)
	uh := userHandler.NewUserHandler(userUC)
	routes.UserRoute(uh, mux)

	// Product
	productRe := productDB.NewProductRepoPG(pg)
	productUC := productUsecase.NewProductUseCase(productRe)
	pH := productHandler.NewProductHandler(productUC)
	routes.ProductRoute(pH, mux)

	// Order
	orderRepo := orderDB.NewOrderRepoPg(pg)
	orderUC := orderUsecase.NewOrderUseCase(orderRepo)
	orderH := orderHandler.NewOrderHandler(orderUC)
	routes.OrderRoute(orderH, mux)

	wrappedMux := middleware.AuthMiddleware(mux)

	// Combine the two ServeMux instances
	http.Handle("/", wrappedMux)
	http.Handle("/auth/", http.StripPrefix("/auth", authMux))

	log.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
