package handlers

import (
	"log"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func HandleSignup(c *fiber.Ctx) error {

	log.Println(c.Body())
	// Store the user name and password in the database
	// Create UID
	// Start a connection to the database
	// Store the info
	// close the connection
	// send the response back
	
	uid := utils.GetNumber()

	db := database.CreateConnection()

	sqlStatement := `
		INSERT INTO Users (id, username, hashed_password)
		VALUES ($1, $2, $3)
		`

	_, err := db.Exec(sqlStatement, uid, "aditya15", "password")
	if err != nil {
		log.Panic(err)
	}
	log.Println("User Created")

	database.CloseDB(db);
	
	// return a response
	c.SendString("OK")

	return nil
}

func HandleLogin(c *fiber.Ctx) error {
	return nil
}