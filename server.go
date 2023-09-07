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
	db.DB.AutoMigrate(&model.User{}, &model.ResetPassword{})
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
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
	})

	// Static folder
	app.Static("/", "./public")

	page := app.Group("/", middleware.IsLoggedIn)
	htmx := app.Group("/")
	//
	// Routes which render sites
	//
	page.Get("/", func(c *fiber.Ctx) error {
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
	page.Get("register", middleware.GuestOnly, controller.RenderRegister)
	page.Get("login", middleware.GuestOnly, controller.RenderLogin)
	page.Get("logout", controller.LogoutUser)
	page.Get("user/verify/:uuid", controller.VerifyUser)
	page.Get("user/update", middleware.LoginRequired, controller.RenderUpdateUser)
	page.Post("user/updatepassword", middleware.LoginRequired, controller.UpdateUserPassword)
	page.Get("auth/resetpassword", controller.RenderResetPassword)
	page.Get("auth/resetpassword/set/:key", controller.RenderResetPasswordSet)
	page.Get("contact", controller.RenderContactFormular)
	//
	// Routes which return HTML Chuncks for HTMX
	//
	htmx.Post("register", controller.CreateUser)
	htmx.Post("login", controller.LoginUser)
	htmx.Post("auth/resetpassword", controller.ResetPassword)
	htmx.Post("auth/resetpassword/set/:key", controller.ResetPasswordSet)
	htmx.Post("user/update", middleware.LoginRequired, controller.UpdateUserProfile)
	htmx.Post("contact", controller.SubmitContactFormular)

	// Start server
	log.Fatal(app.Listen("127.0.0.1:3000"))
}
