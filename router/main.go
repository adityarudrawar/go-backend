package router

import (
	"github.com/adityarudrawar/go-backend/handlers"
	"github.com/gofiber/fiber/v2"
)


func SetupRoutes(app *fiber.App) {
	// Base: Ping
	app.Get("/health", handlers.HandleHealthCheck)

	auth := app.Group("/register")
	auth.Post("/signup", handlers.HandleSignup)
	auth.Post("/login", handlers.HandleLogin)

	message := app.Group("/message")
	message.Post("/", handlers.HandlePostMessage)
	message.Get("/", handlers.HandleGetMessage)
	message.Post("/upvote/:id", handlers.HandleUpvote)
	message.Post("/downvote/:id", handlers.HandleDownvote)
}