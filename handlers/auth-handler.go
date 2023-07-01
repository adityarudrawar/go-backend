package handlers

import (
	"log"
	"time"

	"net/http"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Read the secret from env
// TODO: Abstract getting token and username to a function

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

	var existingUser models.User

	result := database.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		statusCode := fiber.StatusBadRequest
		response := models.BuildResponse(http.StatusText(statusCode), "Username already exists", nil, "Username already exists")
		return c.Status(statusCode).JSON(response)
	}
	
	u := database.DB.Create(&user)

	if u.Error != nil {
		statusCode := fiber.StatusInternalServerError
		response := models.BuildResponse(http.StatusText(statusCode), "Failed to create new user", nil, u.Error.Error())
		return c.Status(statusCode).JSON(response)
	}
	
	statusCode := fiber.StatusCreated
	response := models.BuildResponse(http.StatusText(statusCode), "User Created succesfully", user, "")
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

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		statusCode := fiber.StatusBadRequest
		response := models.BuildResponse(http.StatusText(statusCode), "Incorrect password", nil, "Incorrect password")
		return c.Status(statusCode).JSON(response)
	}

	var session models.Session

	qr := database.DB.Where("username = ?", user.Username).First(&session)
	if qr.Error == nil {
		statusCode := fiber.StatusBadRequest
		response := models.BuildResponse(http.StatusText(statusCode), "User already logged in", nil, "User already logged in")
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

	session = models.Session{
		Username: user.Username,
		Jwt: cookie.Value,
		Expires: cookie.Expires.Unix(),
	}

	result := database.DB.Create(&session)

	if result.Error != nil {
		statusCode := fiber.StatusInternalServerError
		response := models.BuildResponse(http.StatusText(statusCode), "Server Error", user, result.Error.Error())
		return c.Status(statusCode).JSON(response)
	}
	
	statusCode := fiber.StatusOK
	response := models.BuildResponse(http.StatusText(statusCode), "Logged in", user, "")
	return c.Status(statusCode).JSON(response)
}

func HandleLogout(c *fiber.Ctx) error {

	authorized, err := isAuthenticated(c)
	if err != nil {
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	if !authorized {
		return c.Status(fiber.StatusOK).SendString("NOT Authorized")
	}
	
	cookie := c.Cookies("jwt")
	
	cookieToken, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	claims := cookieToken.Claims.(*jwt.StandardClaims)

	var session models.Session

	qr := database.DB.Where("username = ? ", claims.Issuer).Delete(&session)
	if qr.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(qr.Error.Error())
	}

	newCookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&newCookie)

	return c.Status(fiber.StatusOK).SendString("Logged out")
}

func HandleIsAuthenicated(c * fiber.Ctx) error {

	authorized, err := isAuthenticated(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if !authorized {
		return c.Status(fiber.StatusUnauthorized).SendString("Unathorized")
	}

	return c.Status(fiber.StatusOK).SendString("Authorized")
}

func isAuthenticated(c *fiber.Ctx) (bool, error) {
	cookie := c.Cookies("jwt")
	
	cookieToken, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	
	if err != nil {
		return false, nil;
	}

	claims := cookieToken.Claims.(*jwt.StandardClaims)

	var session models.Session

	qr := database.DB.Where("username = ?", claims.Issuer).First(&session)

	if qr.Error != nil {
		return false, qr.Error
	}

	if qr.RowsAffected == 0 {
		return false, nil
	}

	log.Println(cookieToken.Raw)
	// Check the token against the one in database and the the expriry
	if session.Jwt != cookieToken.Raw {
		return false, nil
	}

	storedTime := time.Unix(session.Expires, 0)
	currentTime := time.Now()

	// Compare the current time with the stored time
	if currentTime.After(storedTime) {
		return false, nil
	}

	return true, nil
}