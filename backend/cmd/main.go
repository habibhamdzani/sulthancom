package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sulthancom/pos-backend/internal/config"
	"github.com/sulthancom/pos-backend/internal/database"
	"github.com/sulthancom/pos-backend/internal/handlers"
	"github.com/sulthancom/pos-backend/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
		BodyLimit: 10 * 1024 * 1024,
	})

	middleware.Setup(app)

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redis, err := database.NewRedis(cfg)
	if err != nil {
		log.Printf("Redis connection failed, continuing without cache: %v", err)
	}
	if redis != nil {
		defer redis.Close()
	}

	handlers.Setup(app, db, redis, cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
