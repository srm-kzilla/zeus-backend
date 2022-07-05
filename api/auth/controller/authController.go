package authController

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	authModel "github.com/srm-kzilla/events/api/auth/model"
	authService "github.com/srm-kzilla/events/api/auth/service"
	"github.com/srm-kzilla/events/database"
	"github.com/srm-kzilla/events/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

/***********************
To register a new Admin.
***********************/
func RegisterAdmin(c *fiber.Ctx) error {
	var user authModel.User
	c.BodyParser(&user)

	if errors := validators.ValidateAdminUser(user); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var check authModel.User
	adminCollection, err := database.GetCollection(os.Getenv("DB_NAME"), "Admin")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	adminCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&check)
	if check.Email == user.Email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}
	user.ID = primitive.NewObjectID()
	pwd, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(pwd)
	user.CreatedAt = time.Now()

	res, e := adminCollection.InsertOne(context.Background(), user)
	if e != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":      e.Error(),
			"InsertedId": res.InsertedID,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Admin user created successfully",
		"email":   user.Email,
	})
}

/***********************
To log in the registered admin.
***********************/
func LoginAdmin(c *fiber.Ctx) error {
	var user authModel.User
	c.BodyParser(&user)

	adminCollection, err := database.GetCollection(os.Getenv("DB_NAME"), "Admin")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var check authModel.User
	error := adminCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&check)
	if error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email does not exist",
		})
	}
	e := bcrypt.CompareHashAndPassword([]byte(check.Password), []byte(user.Password))
	if e != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is incorrect",
		})
	}
	token, err := authService.GenerateToken(user.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	refresh, err := authService.GenerateRefreshToken(user.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"refresh": refresh,
		"email":   user.Email,
	})
}

/***********************
New Auth Token generator for an Admin User.
***********************/
func RefreshAdmin(c *fiber.Ctx) error {
	refreshToken := c.GetReqHeaders()["X-Refresh-Token"]
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No refresh token provided",
		})
	}

	token, err := authService.AuthenticateToken(refreshToken, "REFRESH")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid refresh token",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var email string = claims.Issuer

	adminCollection, err := database.GetCollection(os.Getenv("DB_NAME"), "Admin")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var check authModel.User
	error := adminCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&check)
	if error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email does not exist",
		})
	}
	accessToken, err := authService.GenerateToken(email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Refresh successful",
		"token":   accessToken,
		"refresh": refreshToken,
	})
}

/***********************
Authorization Protocol allowing access to only-admin accessible routes.
***********************/
func AuthenticateAdmin(c *fiber.Ctx) error {
	accessToken := c.GetReqHeaders()["X-Access-Token"]
	if accessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No access token provided",
		})
	}
	token, err := authService.AuthenticateToken(accessToken, "SECRET")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid access token",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var email string = claims.Issuer

	adminCollection, err := database.GetCollection(os.Getenv("DB_NAME"), "Admin")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var check authModel.User
	error := adminCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&check)
	if error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email does not exist",
		})
	}
	fmt.Println("Admin : ", email)
	return c.Next()
}
