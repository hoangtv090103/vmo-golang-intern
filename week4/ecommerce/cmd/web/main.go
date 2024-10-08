// main.go
package main

import (
	"database/sql"
	accountHandler "ecommerce/internal/account/handler"
	"ecommerce/internal/account/infra"
	"ecommerce/internal/account/usecase"
	"ecommerce/pkg/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

const port = `:8080`

type application struct {
	accountHandler *accountHandler.AccountHandler
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

func setupApplication(database *sql.DB) *application {
	// Initialize repository
	accountRepo := infra.NewAccountPGRepository(database)

	// Initialize usecase
	accountUsecase := usecase.NewAccountUsecase(accountRepo)

	// Initialize handler
	ah := accountHandler.NewAccountHandler(accountUsecase)

	return &application{
		accountHandler: ah,
	}
}

func (app *application) setupRoutes(fiberApp *fiber.App) {
	// Account routes
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})
	// fiberApp.Post("/register", app.accountHandler.Register)
	fiberApp.Post("/login", app.accountHandler.Login)
	fiberApp.Post("/register", app.accountHandler.Register)
}
