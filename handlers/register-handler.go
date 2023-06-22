package handlers

import (
	"log"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type SignupInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupOutput struct { 
	Username string `json:"username" binding:"required"`
	Id string `json:"id" binding:"required"`
}

type ErrorMessageOutput struct {
	Message string `json:"message"`
}

func HandleSignup(c *fiber.Ctx) error {
	uid := utils.GetNumber()

	// From request body get the json
	signupInput := new(SignupInput)
	if err := c.BodyParser(signupInput); err != nil {
        log.Println("error = ",err)
        return c.SendStatus(200)
    }

	// from json extract username and password
	name := signupInput.Username
	password := signupInput.Password


	// return generated uid.
	db := database.CreateConnection()

	sqlStatement := `
		INSERT INTO Users (id, username, hashed_password)
		VALUES ($1, $2, $3)
		`
	// store that username and password in database
	_, err := db.Exec(sqlStatement, uid, name, password)
	if err != nil {
		log.Panic(err)
		// error: conflicting username
		message := "Request failed, username taken" 	
		errorMessage := ErrorMessageOutput{
			Message: message,
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "not successful",
			"data":   errorMessage,
		})
	}

	database.CloseDB(db);
	
	// return a response with the uid and a message
	signupOutput := SignupOutput{
		Id : uid,	
		Username : name,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   signupOutput,
	})
}



func HandleLogin(c *fiber.Ctx) error {
	// parses username, and return id if the user exists, if not returns invalid request: User not signed up.

	return nil
}