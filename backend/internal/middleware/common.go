package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(app *fiber.App) {
	app.Use(recover.New())

	app.Use(logger.New(logger.Config{
		Format:     "${pid} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Response-Time", time.Now().Format(time.RFC3339Nano))
		return c.Next()
	})
}
