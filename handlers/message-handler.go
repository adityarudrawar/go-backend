package handlers

import (
	"log"
	"time"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type MessageInput struct {
	Sender  string `json:"sender" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type Message struct {
	Id string `json:"id"`
	Sender string `json:"sender"`
	CreatedAt string `json:"createdAt"`
	Content string `json:"content"`
	Upvotes int `json:"upvotes"`
	Downvotes int `json:"downvotes"`
}

func HandlePostMessage(c *fiber.Ctx) error {
	messageInput := new(MessageInput)
	if err := c.BodyParser(messageInput); err != nil {
		log.Println(err)
        return c.SendStatus(200)
    }

	uid := utils.GetNumber(15)
	created_at := time.Now()
	sender := messageInput.Sender
	content := messageInput.Content

	db := database.CreateConnection()

	sqlStatement := `
		INSERT INTO Messages (id, created_at, sender, upvotes, downvotes, content)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(sqlStatement, uid, created_at, sender, 0, 0, content)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data":   "Sending message failed",
		})
	}

	database.CloseDB(db)

	message := Message{
		Id : uid,
		Sender : sender,
		CreatedAt : created_at.String(),
		Content : content,
		Upvotes : 0,
		Downvotes : 0,
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   message,

	})
}


func HandleGetMessage(c *fiber.Ctx) error {
	// If the user is signedin
	db := database.CreateConnection()

	var messages []Message
	rows, err := db.Query("SELECT * FROM Messages ORDER BY created_at DESC")
	if err != nil {
		log.Panic(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data":   "Getting messages failed",
		})
	}

	defer rows.Close()
	for rows.Next() {
		var message Message
		
		if err := rows.Scan(&message.Id, &message.CreatedAt, &message.Sender,&message.Upvotes, &message.Downvotes, &message.Content); err != nil {
			log.Panic(err)
			message := err.Error()
			errorMesage := ErrorMessageOutput {
				Message: message,
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "not successful",
				"data":   errorMesage,
			})
		}
		messages = append(messages, message)
	}
	
	database.CloseDB(db)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   messages,
	})
}

func HandleUpvote(c *fiber.Ctx) error {
	// If the user is signedin
    id := c.Params("id")

    db := database.CreateConnection()
    defer database.CloseDB(db)

    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    _, err = tx.Exec("UPDATE Messages SET upvotes = upvotes + 1 WHERE id = $1", id)
    if err != nil {
        tx.Rollback()
		// TODO: Add the error message object
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status": "error",
            "data":   "Could not upvote",
        })
    }

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }


	sqlStatement := `
		SELECT * FROM Messages WHERE id = $1
	`

	row := db.QueryRow(sqlStatement, id)

	// var id string 
	var createdAt string
	var sender string
	var upvotes int
	var downvotes int
	var content string

	err = row.Scan(&id, &createdAt, &sender, &upvotes, &downvotes, &content)
	if err != nil {
		log.Panic(err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   err.Error(),
		})
	}

	message := Message{
		Id : id,
		CreatedAt: createdAt,
		Sender: sender,
		Upvotes: upvotes,
		Downvotes: downvotes,
		Content: content,
	}

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   message,
    })
}


func HandleDownvote(c *fiber.Ctx) error {
	// If the user is signedin

    id := c.Params("id")

    db := database.CreateConnection()
    defer database.CloseDB(db)

    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    _, err = tx.Exec("UPDATE Messages SET downvotes = downvotes + 1 WHERE id = $1", id)
    if err != nil {
        tx.Rollback()
		// TODO: Add the error message object
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status": "error",
            "data":   "Could not downvote",
        })
    }

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }


	sqlStatement := `
		SELECT * FROM Messages WHERE id = $1
	`

	row := db.QueryRow(sqlStatement, id)

	// var id string 
	var createdAt string
	var sender string
	var upvotes int
	var downvotes int
	var content string

	err = row.Scan(&id, &createdAt, &sender, &upvotes, &downvotes, &content)
	if err != nil {
		log.Panic()
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   err.Error(),
		})
	}

	message := Message{
		Id : id,
		CreatedAt: createdAt,
		Sender: sender,
		Upvotes: upvotes,
		Downvotes: downvotes,
		Content: content,
	}

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   message,
    })
}
