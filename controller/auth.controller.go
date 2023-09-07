package controller

import (
	"time"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/nickfthedev/fiberHTMX/db"
	"github.com/nickfthedev/fiberHTMX/lib"
	"github.com/nickfthedev/fiberHTMX/model"
	"golang.org/x/crypto/bcrypt"
)

// PAGE | Renders Reset Password Page where you set the password
func RenderResetPasswordSet(c *fiber.Ctx) error {
	key := c.Params("key") // Secret Key
	// Check DB for Key
	//...
	rp := new(model.ResetPassword)
	db.DB.Where("key = ?", key).First(&rp)
	// Create database entry
	if rp.ID == 0 {
		return c.Redirect("/")
	}
	// Token expired
	if time.Now().After(rp.UpdatedAt.Add(60 * time.Minute)) {
		return c.Render("auth/resetpassword", fiber.Map{"ErrorMessage": "Your Key is expired. Please request a new email."})
	}
	return c.Render("auth/resetpasswordset", fiber.Map{"key": key})
}

// HTMX | Save the new password
func ResetPasswordSet(c *fiber.Ctx) error {
	key := c.Params("key") // Secret Key
	newPassword := c.FormValue("password")
	confirmPassword := c.FormValue("confirmpassword")
	//Check if Password and ConfirmPassword Match
	if newPassword != confirmPassword {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Passwords does not match"}, "common/empty")
		//return c.Status(500).JSON(fiber.Map{"error": "Passwords does not match"})
	}
	//Check if password has at least 6 characters, one uppercase, one lowercase and one number
	password := newPassword
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
	}

	// Check for the key
	rp := new(model.ResetPassword)
	db.DB.Where("key = ?", key).First(&rp)
	// Create database entry
	if rp.ID == 0 {
		return c.Redirect("/")
	}

	// Get user from DB
	user := new(model.User)
	db.DB.Where("uuid = ?", rp.UserUUID).First(&user)
	if user.ID == 0 {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Unexpected Error. User not found. Please contact us for resetting your password."}, "common/empty")
	}
	// Change password from user to new password
	//Hash Password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to hash password"}, "common/empty")
	}
	user.Password = string(hash)
	db.DB.Save(&user)
	// Delete entry from ResetPassword Table
	db.DB.Unscoped().Delete(&rp)
	c.Append("HX-Trigger", "myEvent") // Send to redirect on success
	return c.Render("common/success", fiber.Map{"SuccessMessage": "Password has been changed successfully"}, "common/empty")
}

// PAGE | Renders Reset Password Page
func RenderResetPassword(c *fiber.Ctx) error {
	return c.Render("auth/resetpassword", fiber.Map{})
}

// PAGE | Send a mail with a link for setting new password
func ResetPassword(c *fiber.Ctx) error {
	email := c.FormValue("email")

	// Find the user by mail or throw an error
	user := new(model.User)
	db.DB.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Unknown User"}, "common/empty")
	}
	// Try to find Entry in DB
	rp := new(model.ResetPassword)
	db.DB.Where("user_uuid = ?", user.UUID).First(&rp)
	// Create database entry
	if rp.ID == 0 { // No entry for now create a new one
		rp.UserUUID = user.UUID
		rp.Key = uuid.New()
		db.DB.Create(&rp)
	} else { // Update the existing entry for password reset
		rp.Key = uuid.New()
		db.DB.Save(&rp)
	}
	// Send mail with link
	msg := "Here is your link for resetting the password for " + lib.Config.AppName + "!+<br><br> Please click the link below to change your password: <br><br> <a href=\"https://" + lib.Config.Host + "/auth/resetpassword/set/" + rp.Key.String() + "\">Reset Password</a><br><br><i>The link is only valid for 1 hour</i>"
	errMail := lib.SendEmail(user.Email, "Reset your password | "+lib.Config.AppName, msg, "text/html")
	if errMail != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to send Mail. Please try again"}, "common/empty")
	}
	// Success, no error
	return c.Render("common/success", fiber.Map{"SuccessMessage": "A Mail has been sent to your email account. Click the link in the email to reset your password"}, "common/empty")
}

// PAGE | Render register
func RenderRegister(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"Title": "Hello, World!",
	})
}

// PAGE | Render login
func RenderLogin(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{})
}

// PAGE (Only a redirect) | Logout User & Destroy Cookie
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

// PAGE
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

	if !user.Verified {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "You're not verified yet", "ErrorCode": "Please check your email inbox and spam folder"}, "common/empty")
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
