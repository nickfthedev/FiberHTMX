package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nickfthedev/fiberHTMX/lib"
)

// PAGE | Renders Reset Password Page where you set the password
func RenderContactFormular(c *fiber.Ctx) error {
	return c.Render("contact/formular", fiber.Map{})
}

// HTMX | Submit Contact Formular
func SubmitContactFormular(c *fiber.Ctx) error {
	// HTMX | Submit Contact Formular
	email := c.FormValue("email")
	name := c.FormValue("name")
	message := c.FormValue("message")

	if email == "" || name == "" || message == "" {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Please fill out all fields"}, "common/empty")
	}

	// Send email to lib.Config.SMTPUser
	// Send mail with link
	msg := "New Contact Request from " + name + " for " + lib.Config.AppName + "!+<br><br> Here is the message: <br><br>" + message
	errMail := lib.SendEmail(lib.Config.SMTPUser, "Contact Formular | "+lib.Config.AppName, msg, "text/html")
	if errMail != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to send Mail. Please try again"}, "common/empty")
	}

	return c.Render("common/success", fiber.Map{"SuccessMessage": "Message sent successfully"}, "common/empty")
}
