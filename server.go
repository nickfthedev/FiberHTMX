package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/nickfthedev/fiberHTMX/controller"
	"github.com/nickfthedev/fiberHTMX/db"
	"github.com/nickfthedev/fiberHTMX/lib"
	"github.com/nickfthedev/fiberHTMX/model"

	"gorm.io/gorm"
)

type User struct {
	IsLoggedIn bool
	// Other user-related data
}

func main() {
	//
	//Config
	//
	lib.LoadConfig("config.json")
	fmt.Println("Config Loaded:", lib.Config.DbDriver)

	//
	// Database
	//
	db.ConnectDB(lib.Config)
	db.DB.AutoMigrate(&model.User{})
	// Create standard user if no user is found in database
	// TODO Firstorcreate GORM
	var u model.User
	r := db.DB.Take(&u)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		// Create Standard User: admin Password: password
		controller.CreateStandardAdminUser()
	}
	// Simulate user authentication status (you should implement your actual login logic)
	user := User{IsLoggedIn: false}

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Static folder
	app.Static("/", "./public")

	//
	// Routes
	//
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"IsLoggedIn": user.IsLoggedIn,
			"Title":      "Hello, World!",
		})
	})
	app.Get("/register", controller.RenderRegister)
	app.Get("/login", controller.RenderLogin)
	app.Post("register", controller.CreateUser)

	// Start server
	log.Fatal(app.Listen("127.0.0.1:3000"))
}
