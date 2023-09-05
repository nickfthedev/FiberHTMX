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
	"github.com/nickfthedev/fiberHTMX/middleware"
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

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Static folder
	app.Static("/", "./public")

	page := app.Group("/", middleware.IsLoggedIn)
	htmx := app.Group("/")
	//
	// Routes which render sites
	//
	page.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Locals("id"))
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})
	page.Get("/protected", middleware.LoginRequired, func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Protected Page!",
		})
	})
	page.Get("/register", controller.RenderRegister)
	page.Get("/login", controller.RenderLogin)
	page.Get("logout", controller.LogoutUser)

	//
	// Routes which return HTML Chuncks for HTMX
	//
	htmx.Post("register", controller.CreateUser)
	htmx.Post("login", controller.LoginUser)

	// Start server
	log.Fatal(app.Listen("127.0.0.1:3000"))
}
