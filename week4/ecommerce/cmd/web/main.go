// main.go
package main

import (
	accountHandler "ecommerce/internal/auth/handler"
	orderHandler "ecommerce/internal/order/handler"
	productHandler "ecommerce/internal/product/handler"
	"ecommerce/internal/user/userHandler"

	"ecommerce/pkg/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

const port = `:8080`

type application struct {
	accountHandler *accountHandler.AccountHandler
	userHandler    *userHandler.UserHandler
	productHandler *productHandler.ProductHandler
	orderHandler   *orderHandler.OrderHandler
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbInstance := db.GetDBInstance()

	defer dbInstance.Close()

	// Initialize application
	app := setupApplication(dbInstance)

	fiberApp := fiber.New()

	// Setup routes
	app.setupRoutes(fiberApp)

	log.Fatal(fiberApp.Listen(port))
}

func (app *application) setupRoutes(fiberApp *fiber.App) {
	// Account routes
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	// Auth routes
	fiberApp.Post("/login", app.accountHandler.Login)
	fiberApp.Post("/register", app.accountHandler.Register)

	// User routes
	fiberApp.Get("/users", app.userHandler.GetAllUsers)
	fiberApp.Post("/users", app.userHandler.AddUser)
	fiberApp.Put("/users/:id", app.userHandler.UpdateUser)
	fiberApp.Delete("/users/:id", app.userHandler.DeleteUser)

	// Product routes
	fiberApp.Get("/products", app.productHandler.GetAllProducts)
	fiberApp.Get("/products/:id", app.productHandler.GetProductByID)
	fiberApp.Post("/products", app.productHandler.AddProduct)
	fiberApp.Put("/products/:id", app.productHandler.UpdateProduct)
	fiberApp.Delete("/products/:id", app.productHandler.DeleteProduct)

	// Order routes
	fiberApp.Get("/orders", app.orderHandler.GetAllOrders)
	fiberApp.Get("/orders/:username", app.orderHandler.GetUserOrders)
	fiberApp.Post("/orders", app.orderHandler.CreateOrder)
	fiberApp.Put("/orders/:id", app.orderHandler.UpdateOrder)
	fiberApp.Delete("/orders/:id", app.orderHandler.DeleteOrder)
	fiberApp.Get("/orders/:id/invoice", app.orderHandler.GetInvoice)
	fiberApp.Get("/orders/:id/print-invoice", app.orderHandler.PrintInvoice)

}
