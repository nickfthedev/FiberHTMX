package controller

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/nickfthedev/fiberHTMX/db"
	"github.com/nickfthedev/fiberHTMX/model"
	"golang.org/x/crypto/bcrypt"
)

// func CreateUser
func CreateUser(c *fiber.Ctx) error {
	//Check if all fields are filled
	checkUser := new(model.RegisterUser)
	if err := c.BodyParser(checkUser); err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Review your input!"}, "common/empty")
		//return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}
	//Check if Password and ConfirmPassword Match
	if checkUser.Password != checkUser.ConfirmPassword {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Passwords does not match"}, "common/empty")
		//return c.Status(500).JSON(fiber.Map{"error": "Passwords does not match"})
	}
	//Check if password has at least 6 characters, one uppercase, one lowercase and one number
	password := checkUser.Password
	if len(password) < 6 {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Password must be at least 6 characters long"}, "common/empty")
		//return c.Status(500).JSON(fiber.Map{"error": "Password must be at least 6 characters long"})
	}
	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Password must have at least one uppercase, one lowercase and one number"}, "common/empty")
		//return c.Status(500).JSON(fiber.Map{"error": "Password must have at least one uppercase, one lowercase and one number"})
	}

	//Map input to user Model
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Review your input"}, "common/empty")
		//return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	// Generate a random user name from name field
	randString := func(n int) string {
		letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		s := make([]rune, n)
		for i := range s {
			s[i] = letters[rand.Intn(len(letters))]
		}
		return string(s)
	}

	user.Username = user.Name + randString(5) // append 5 random characters to the name

	//Hash Password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to hash password"}, "common/empty")
		//return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to hash password", "data": err.Error()})
	}
	user.Password = string(hash)
	result := db.DB.Create(&user) // pass pointer of data to Create
	if result.Error != nil {      //Error handling
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
			return c.Render("common/error", fiber.Map{"ErrorMessage": "Registration failed", "ErrorCode": "E-Mail is already in use "}, "common/empty")
		}
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to create user", "ErrorCode": result.Error.Error()}, "common/empty")
	}
	//If no error send back ok
	return c.Render("common/success", fiber.Map{"SuccessMessage": "Registered successfully", "SuccessCode": "You can now login"}, "common/empty")
	//return c.JSON(fiber.Map{"status": "success", "message": "User successfully created", "data": user.ID})

} // func CreateUser(c *fiber.Ctx)

func CreateStandardAdminUser() {
	//Map input to user Model
	user := new(model.User)
	user.Username = "admin"
	user.Name = "admin"
	user.Fullname = "Administator"
	user.Email = "admin@admin.com"
	user.Password = "password"
	user.Verified = true
	//Hash Password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)
	db.DB.Create(&user) // pass pointer of data to Create
	fmt.Println("\n=============================\n", "No User in database. Created user 'admin@admin.com' with password 'password'!", "\n=============================")
}
