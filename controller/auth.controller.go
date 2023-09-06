package controller

import (
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
		"Title": "Hello, World!",
	})
}

func RenderLogin(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{})
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
	return c.Redirect("/", 302)

}

// func LoginUser
func LoginUser(c *fiber.Ctx) error {

	var input model.UserLoginInput // Validate input
	if err := c.BodyParser(&input); err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Review your input!"}, "common/empty")
	}

	// Find user if exists
	var user model.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Email or Password incorrect"}, "common/empty")
	}
	//Check Password against Database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Email or Password incorrect"}, "common/empty")

	}

	//Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), //Expires after 30 Days
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(lib.Config.TokenSecret))
	if err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to create token"}, "common/empty")
	}

	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(30 * 24 * time.Hour)
	cookie.SameSite = "lax"

	// Set cookie
	c.Cookie(cookie)

	//Return Token in Response
	c.Append("HX-Trigger", "myEvent")
	return c.Render("common/success", fiber.Map{"SuccessMessage": "Logged in successfully"}, "common/empty")
}
