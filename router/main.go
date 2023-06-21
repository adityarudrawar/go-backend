package router

import (
	"github.com/adityarudrawar/go-backend/handlers"
	"github.com/gofiber/fiber/v2"
)
func SetupRoutes(app *fiber.App) {
	app.Get("/health", handlers.HandleHealthCheck)
}