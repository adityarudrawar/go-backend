package handlers

import (
	"log"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/utils"
	"github.com/gofiber/fiber/v2"

	"database/sql"
)

type SignupInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupOutput struct { 
	Id string `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type ErrorMessageOutput struct {
	Message string `json:"message"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginOutput struct {
	Username string `json:"username" binding:"required"`
	Id string `json:"id" binding:"required"`
}

func HandleSignup(c *fiber.Ctx) error {
	uid := utils.GetNumber()

	// From request body get the json
	signupInput := new(SignupInput)
	if err := c.BodyParser(signupInput); err != nil {
		// TODO: correct the response here  
        log.Println("error = ",err)
		message := err.Error()	
		errorMessage := ErrorMessageOutput{
			Message: message,
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "not successful",
			"data":   errorMessage,
		})		
    }

	// from json extract username and password
	name := signupInput.Username
	password := signupInput.Password

	db := database.CreateConnection()

	sqlStatement := `
		INSERT INTO Users (id, username, hashed_password)
		VALUES ($1, $2, $3)
		`
	// store that username and password in database
	_, err := db.Exec(sqlStatement, uid, name, password)
	if err != nil {
		log.Panic(err)
		// error: username already exists
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

	// From request body get the json
	loginInput := new(LoginInput)
	if err := c.BodyParser(loginInput); err != nil {
        log.Println("error = ",err)
        return c.SendStatus(200)
    }

	// from json extract username and password
	name := loginInput.Username
	inputPassword := loginInput.Password

	db := database.CreateConnection()

	sqlStatement := `
		SELECT * FROM Users WHERE username = $1
		`
	// store that username and password in database
	row := db.QueryRow(sqlStatement, name)
	
	var id string
	var username string
	var password string

	database.CloseDB(db);
	
	err := row.Scan(&id, &username, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			message := "No user found with this username"
			errorMessage := ErrorMessageOutput{
				Message: message,
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "not successful",
				"data":   errorMessage,
			})
		} else {
			log.Fatal(err)
			message := string(err.Error())
			errorMessage := ErrorMessageOutput{
				Message: message,
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "not successful",
				"data":   errorMessage,
			})
		}
	} else {
		// Password does not match
		if password != inputPassword {
			message := "Incorrect Password"
			errorMessage := ErrorMessageOutput{
				Message: message,
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "not successful",
				"data":   errorMessage,
			})
		}
	}

	loginOutput := LoginOutput{
		Id : id,
		Username : username,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   loginOutput,
	})
}