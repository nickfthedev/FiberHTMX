package main

import (
    "log"
		"fmt"
    "github.com/gofiber/fiber/v2"

		"github.com/nickfthedev/fiberHTMX/lib"
)

func main() {
	
		lib.LoadConfig("config.json");
		fmt.Println("Config:", lib.Cfg.DbDriver);

    app := fiber.New()

    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    log.Fatal(app.Listen(":3000"))
}