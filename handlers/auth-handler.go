package handlers

import (
	"os"
	"strings"
	"time"

	"net/http"

	"github.com/adityarudrawar/go-backend/database"
	"github.com/adityarudrawar/go-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Abstract getting token and username to a function
// TODO: Better responses


var SecretKey = os.Getenv("TOKEN_SECRET")

func HandleRegister(c *fiber.Ctx) error {
	
	var data map[string]string
	
	if err := c.BodyParser(&data); err != nil {
		return err
	} 

	// TODO: Change encryption algorithm
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

	var session models.UserSessions

	// TODO: If a users entry exists in session table but the token has expired, and the user puts a request for login the session should update,
	qr := database.DB.Where("username = ?", user.Username).First(&session)
	if qr.Error == nil {
		statusCode := fiber.StatusBadRequest
		response := models.BuildResponse(http.StatusText(statusCode), "User already logged in", nil, "User already logged in")
		return c.Status(statusCode).JSON(response)
	}
	
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 Hour
	})
	
	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		response := models.BuildResponse(http.StatusText(statusCode), "Could not login", nil, "Error signing string")
		return c.Status(statusCode).JSON(response)
	}

	session = models.UserSessions{
		Username: user.Username,
		Token: token,
	}

	result := database.DB.Create(&session)

	if result.Error != nil {
		statusCode := fiber.StatusInternalServerError
		response := models.BuildResponse(http.StatusText(statusCode), "Error creating session", user, result.Error.Error())
		return c.Status(statusCode).JSON(response)
	}
	
	c.Set("Authorization", "Bearer "+token)

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

	requestToken := c.Get("Authorization")
	
	requestToken = strings.TrimPrefix(requestToken, "Bearer ")
	
	_, claims, err := getJwtTokenAndClaims(requestToken)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		response := models.BuildResponse(http.StatusText(statusCode), "Error Logging out", nil, err.Error())
		return c.Status(statusCode).JSON(response)
	}

	var session models.UserSessions

	qr := database.DB.Where("username = ? ", claims.Issuer).Delete(&session)
	if qr.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(qr.Error.Error())
	}

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

	requestToken := c.Get("Authorization")
	
	requestToken = strings.TrimPrefix(requestToken, "Bearer ")

	jwtToken, claims, err := getJwtTokenAndClaims(requestToken)
	if err != nil {
		return false, err 
	}

	var userSession models.UserSessions

	qr := database.DB.Where("username = ?", claims.Issuer).First(&userSession)

	if qr.Error != nil {
		return false, qr.Error
	}

	if qr.RowsAffected == 0 {
		return false, nil
	}

	if userSession.Token != jwtToken.Raw {
		return false, nil
	}

	storedTime := time.Unix(claims.ExpiresAt, 0)
	currentTime := time.Now()

	// Compare the current time with the stored time
	if currentTime.After(storedTime) {
		return false, nil
	}

	return true, nil
}

func getJwtTokenAndClaims(token string) (*jwt.Token, *jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	
	if err != nil {
		return nil, nil, err
	}

	claims := jwtToken.Claims.(*jwt.StandardClaims)

	return jwtToken, claims, nil;
}