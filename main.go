package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/sqlite3"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	storage := sqlite3.New()
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        250,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		LimiterMiddleware: limiter.FixedWindow{},
		Storage:           storage,
	}))
	app.Get("/fw", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! - From Fix Window")
	})
	app.Get("/fw2", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! - From Fix Window 2")
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error Run Fiber: %s", err.Error())
	}
}
