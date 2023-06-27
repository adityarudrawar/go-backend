package handlers

import (
	"time"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Write better response objects for each error

var SecretKey = "dhkjgbfkljkljbdlkjbfjkb"

func HandleRegister(c *fiber.Ctx) error {
	
	var data map[string]string
	
	if err := c.BodyParser(&data); err != nil {
		return err
	} 

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Username: data["username"],
		Password: password,
	}
	
	u := database.DB.FirstOrCreate(&user)

	if u.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "Internal server error",
			})
	}

	if u.RowsAffected == 1 {
		return c.Status(fiber.StatusCreated).JSON(
			fiber.Map{
				"message": "User created successfully",
			})
	}
	
	return c.Status(fiber.StatusBadRequest).JSON(
		fiber.Map{
			"error": "Username already exists",
		})
}


func HandleLogin(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	dbc := database.DB.Where("username = ?", data["username"]).First(&user)

	if dbc.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Compare the hash of the password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})
	
	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func HandleGetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("username = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func HandleLogout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}