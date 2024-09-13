package handlers

import (
	"time"

	"net/http"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func HandleGetMessages(c *fiber.Ctx) error {
	authorized, err := IsAuthenticated(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if !authorized {
		return c.Status(fiber.StatusUnauthorized).SendString("Unathorized")
	}

	var messages []models.Message
	database.DB.Find(&messages)

	statusCode := fiber.StatusOK
	response := models.BuildResponse(http.StatusText(statusCode), "All Messages", messages, "")
	return c.Status(statusCode).JSON(response)
}

func HandlePostMessage(c * fiber.Ctx) error {

	authorized, err := IsAuthenticated(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if !authorized {
		return c.Status(fiber.StatusUnauthorized).SendString("Unathorized")
	}

	var postMessage models.PostMessageRequest

	if err := c.BodyParser(&postMessage); err != nil {
		// TODO: Build error response
        return err
    }

	var message = models.Message{
		Uuid : uuid.New().String(),
		CreatedAt: time.Now(),
		Sender: postMessage.Username,
		Upvotes: 0,
		Downvotes: 0,
		Content: postMessage.Content,
	}

	m := database.DB.Create(&message)

	if m.Error != nil {
		statusCode := fiber.StatusInternalServerError
		response := models.BuildResponse(http.StatusText(statusCode), "Failed to store the message", nil, m.Error.Error())
		return c.Status(statusCode).JSON(response)
	}

	var messages []models.Message
	database.DB.Order("created_at").Find(&messages)
	
	statusCode := fiber.StatusCreated
	response := models.BuildResponse(http.StatusText(statusCode), "Message posted succesfully", messages, "")
	return c.Status(statusCode).JSON(response)
}

func HandleUpvote(c * fiber.Ctx) error {
	
	authorized, err := IsAuthenticated(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if !authorized {
		return c.Status(fiber.StatusUnauthorized).SendString("Unathorized")
	}

	uuid := c.Params("messageId")

	var message models.Message
	if err := database.DB.Where("uuid = ?", uuid).First(&message).Error; err != nil {
        panic(err)
    }
	
	database.DB.Model(&models.Message{}).Where("uuid = ?", uuid).Update("upvotes", message.Upvotes + 1)

	var messages []models.Message
	database.DB.Order("created_at").Find(&messages)

	
	statusCode := fiber.StatusCreated
	response := models.BuildResponse(http.StatusText(statusCode), "Message posted succesfully", messages, "")
	return c.Status(statusCode).JSON(response)
}

func HandleDownvote(c * fiber.Ctx) error {
		
	authorized, err := IsAuthenticated(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if !authorized {
		return c.Status(fiber.StatusUnauthorized).SendString("Unathorized")
	}

	uuid := c.Params("messageId")

	var message models.Message
	if err := database.DB.Where("uuid = ?", uuid).First(&message).Error; err != nil {
        panic(err)
    }
	
	database.DB.Model(&models.Message{}).Where("uuid = ?", uuid).Update("downvotes", message.Downvotes + 1)

	var messages []models.Message
	database.DB.Order("created_at").Find(&messages)

	
	statusCode := fiber.StatusCreated
	response := models.BuildResponse(http.StatusText(statusCode), "Message posted succesfully", messages, "")
	return c.Status(statusCode).JSON(response)
}
