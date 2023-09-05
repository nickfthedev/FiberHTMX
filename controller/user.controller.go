package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nickfthedev/fiberHTMX/db"
	"github.com/nickfthedev/fiberHTMX/model"
	"golang.org/x/crypto/bcrypt"
)

// func CreateUser
func CreateUser(c *fiber.Ctx) error {
	//Check if all fields are filled
	checkUser := new(model.CreateUser)
	if err := c.BodyParser(checkUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	//Map input to user Model
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err.Error()})
	}

	//Hash Password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to hash password", "data": err.Error()})
	}
	user.Password = string(hash)

	result := db.DB.Create(&user) // pass pointer of data to Create
	if result.Error != nil {      //Error handling
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create user", "data": result.Error.Error()})
	}

	//If no error send back ok
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully created", "data": user.ID})

} // func CreateUser(c *fiber.Ctx)

func CreateStandardAdminUser() {
	//Map input to user Model
	user := new(model.User)
	user.Username = "admin"
	user.Fullname = "Administator"
	user.Email = "admin@admin.com"
	user.Password = "password"
	//Hash Password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)
	db.DB.Create(&user) // pass pointer of data to Create
	fmt.Println("\n=============================\n", "No User in database. Created user 'admin@admin.com' with password 'password'!", "\n=============================")
}
