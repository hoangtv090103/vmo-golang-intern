package main

import (
	"database/sql"
	accountHandler "ecommerce/internal/auth/handler"
	"ecommerce/internal/auth/infra"
	"ecommerce/internal/auth/usecase"
	orderHandler "ecommerce/internal/order/handler"
	orderRepo "ecommerce/internal/order/infra"
	orderUsecase "ecommerce/internal/order/usecase"
	productHandler "ecommerce/internal/product/handler"
	productPGRepo "ecommerce/internal/product/infra"
	productUsecase "ecommerce/internal/product/usecase"
	userInfra "ecommerce/internal/user/infra"
	userUC "ecommerce/internal/user/usecase"
	"ecommerce/internal/user/userHandler"
)

func setupApplication(database *sql.DB) *application {
	// Initialize repository
	accountRepo := infra.NewAccountPGRepository(database)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)
	ah := accountHandler.NewAccountHandler(accountUsecase)

	ur := userInfra.NewUserPGRepository(database)
	uu := userUC.NewUserUsecase(ur)
	uh := userHandler.NewUserHandler(*uu)

	pr := productPGRepo.NewProductPGRepository(database)
	pu := productUsecase.NewProductUsecase(pr)
	ph := productHandler.NewProductHandler(*pu)

	or := orderRepo.NewOrderPGRepository(database)
	ou := orderUsecase.NewOrderUsecase(or)
	oh := orderHandler.NewOrderHandler(ou)

	return &application{
		accountHandler: ah,
		userHandler:    uh,
		productHandler: ph,
		orderHandler:   oh,
	}
}
