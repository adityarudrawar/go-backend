package router

import (
	"github.com/adityarudrawar/go-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/health", handlers.HandleHealthCheck)

	v1.Post("/register", handlers.HandleRegister)
	v1.Post("/login", handlers.HandleLogin)
	v1.Post("/logout", handlers.HandleLogout)

	v1.Get("/authentication", handlers.HandleIsAuthenicated)
	
	v1.Get("/messages", handlers.HandleGetMessages)
	v1.Post("/message", handlers.HandlePostMessage)
	v1.Post("/message/:messageId/upvote", handlers.HandleUpvote)
	v1.Post("/message/:messageId/downvote", handlers.HandleDownvote)
}