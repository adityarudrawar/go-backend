package app

import (
	"os"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/joho/godotenv"
)
func SetupAndRunApp() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// Check Connection to the database
	database.Connect()

	// create app
	app := fiber.New()

	app.Use(cors.New())

	// attach middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	// setup routes
	router.SetupRoutes(app)

	// gets the port and start
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	
	app.Listen(host + ":" + port)

	return nil
}