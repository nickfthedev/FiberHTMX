package main

import (
    "log"

    "github.com/gofiber/fiber/v2"

		"github.com/nickfthedev/fiberHTMX/lib"
)

func main() {
	
		lib.LoadConfig("config.json");

    app := fiber.New()

    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    log.Fatal(app.Listen(":3000"))
}