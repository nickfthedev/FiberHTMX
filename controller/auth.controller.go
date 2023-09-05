package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nickfthedev/fiberHTMX/db"
	"github.com/nickfthedev/fiberHTMX/lib"
	"github.com/nickfthedev/fiberHTMX/model"
	"golang.org/x/crypto/bcrypt"
)

func RenderRegister(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"IsLoggedIn": true,
		"Title":      "Hello, World!",
	})
}

func RenderLogin(c *fiber.Ctx) error {
		"IsLoggedIn": true,
	})
}

// Logout USer & Destroy Cookie
func LogoutUser(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name: "Authorization",
		// Set expiry date to the past
		Expires:  time.Now().Add(-(time.Hour * 2)),
		SameSite: "lax",
	})
	c.Cookie(&fiber.Cookie{
		Name: "UserID",
		// Set expiry date to the past
		Expires:  time.Now().Add(-(time.Hour * 2)),
		SameSite: "lax",
	})
	return c.JSON(fiber.Map{"status": "success", "message": "You've been logged out", "data": nil})

}

// func LoginUser
func LoginUser(c *fiber.Ctx) error {

	var input model.UserLoginInput // Validate input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// Find user if exists
	var user model.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "E-Mail or Password incorrect", "data": err})

	}
	//Check Password against Database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "E-Mail or Password incorrect", "data": err})

	}

	//Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), //Expires after 30 Days
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(lib.Config.TokenSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create token", "data": err})

	}

	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(30 * 24 * time.Hour)
	cookie.SameSite = "lax"

	// Set cookie
	c.Cookie(cookie)
	// ...

	// Create second Cookie with UserID
	cookieUserID := new(fiber.Cookie)
	cookieUserID.Name = "UserID"
	cookieUserID.Value = fmt.Sprint(user.ID)
	cookieUserID.Expires = time.Now().Add(30 * 24 * time.Hour)
	cookieUserID.SameSite = "lax"
	c.Cookie(cookieUserID) //Set Cookie with UserID

	//Return Token in Response
	return c.JSON(fiber.Map{"status": "success", "message": "Successfully loggedin", "data": nil})
} //LoginUser
