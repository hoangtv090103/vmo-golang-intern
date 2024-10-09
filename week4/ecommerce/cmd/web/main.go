// main.go
package main

import (
	accountHandler "ecommerce/internal/auth/handler"
	orderHandler "ecommerce/internal/order/handler"
	productHandler "ecommerce/internal/product/handler"
	"ecommerce/internal/user/userHandler"
	"ecommerce/pkg/middleware"

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

	// Public routes (no auth required)
	fiberApp.Post("/login", app.accountHandler.Login)
	fiberApp.Post("/register", app.accountHandler.Register)

	// Apply auth middleware to all other routes
	api := fiberApp.Group("/api", middleware.AuthMiddleware())

	// Auth routes
	api.Post("/login", app.accountHandler.Login)
	api.Post("/register", app.accountHandler.Register)

	// User routes
	api.Get("/users", app.userHandler.GetAllUsers)
	api.Post("/users", app.userHandler.AddUser)
	api.Put("/users/:id", app.userHandler.UpdateUser)
	api.Delete("/users/:id", app.userHandler.DeleteUser)

	// Product routes
	api.Get("/products", app.productHandler.GetAllProducts)
	api.Get("/products/:id", app.productHandler.GetProductByID)
	api.Post("/products", middleware.IsAdminMiddleware(), app.productHandler.AddProduct)
	api.Put("/products/:id", middleware.IsAdminMiddleware(), app.productHandler.UpdateProduct)
	api.Delete("/products/:id", middleware.IsAdminMiddleware(), app.productHandler.DeleteProduct)

	// Order routes
	api.Get("/orders", middleware.IsAdminMiddleware(), app.orderHandler.GetAllOrders)
	api.Get("/orders/:username", app.orderHandler.GetUserOrders)
	api.Post("/orders", middleware.IsUserMiddleware(), app.orderHandler.CreateOrder)
	api.Put("/orders/:id", app.orderHandler.UpdateOrder)
	api.Delete("/orders/:id", middleware.IsAdminMiddleware(), app.orderHandler.DeleteOrder)
	api.Get("/orders/:id/invoice", app.orderHandler.GetInvoice)
	api.Get("/orders/:id/print-invoice", app.orderHandler.PrintInvoice)

}
