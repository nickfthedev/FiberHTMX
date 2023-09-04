package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/nickfthedev/fiberHTMX/lib"
)

type User struct {
	IsLoggedIn bool
	// Other user-related data
}

func main() {
	//Config
	lib.LoadConfig("config.json")
	fmt.Println("Config Loaded:", lib.Config.DbDriver)

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	// Simulate user authentication status (you should implement your actual login logic)
	user := User{IsLoggedIn: false}

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"IsLoggedIn": user.IsLoggedIn,
			"Title":      "Hello, World!",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
