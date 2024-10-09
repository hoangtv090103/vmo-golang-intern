package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func routes(fa *fiber.App) *fiber.App {
	fa.Use(recover.New())
	fa.Use(timeout.NewWithContext(func(c *fiber.Ctx) error {
		return c.SendString("Request time out")
	}, 60*time.Second))
	return fa
}
