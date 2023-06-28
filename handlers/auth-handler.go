package handlers

import (
	"time"

	"net/http"

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
		statusCode := fiber.StatusInternalServerError
		response := models.BuildResponse(http.StatusText(statusCode), "Failed to create new user", nil, u.Error.Error())
		return c.Status(statusCode).JSON(response)
	}

	if u.RowsAffected == 1 {
		statusCode := fiber.StatusCreated
		response := models.BuildResponse(http.StatusText(statusCode), "User Created succesfully", user, "")
		return c.Status(statusCode).JSON(response)
	}
	
	statusCode := fiber.StatusBadRequest
	response := models.BuildResponse(http.StatusText(statusCode), "Username already exists", nil, "Username already exists")
	return c.Status(statusCode).JSON(response)
}


func HandleLogin(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	dbc := database.DB.Where("username = ?", data["username"]).First(&user)

	if dbc.RowsAffected == 0 {
		statusCode := fiber.StatusNotFound
		response := models.BuildResponse(http.StatusText(statusCode), "User not found", nil, "User not found")
		return c.Status(statusCode).JSON(response)
	}

	// Compare the hash of the password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		statusCode := fiber.StatusBadRequest
		response := models.BuildResponse(http.StatusText(statusCode), "Incorrect password", nil, "Incorrect password")
		return c.Status(statusCode).JSON(response)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})
	
	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		response := models.BuildResponse(http.StatusText(statusCode), "Could not login", nil, "Could not login")
		return c.Status(statusCode).JSON(response)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	
	statusCode := fiber.StatusOK
	response := models.BuildResponse(http.StatusText(statusCode), "Logged in", user, "")
	return c.Status(statusCode).JSON(response)
}

func HandleGetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		statusCode := fiber.StatusUnauthorized
		response := models.BuildResponse(http.StatusText(statusCode), "Unauthenticated", nil, "")
		return c.Status(statusCode).JSON(response)
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("username = ?", claims.Issuer).First(&user)

	statusCode := fiber.StatusOK
	response := models.BuildResponse(http.StatusText(statusCode), "", user, "")
	return c.Status(statusCode).JSON(response)
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