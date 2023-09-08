package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nickfthedev/fiberHTMX/db"
	"github.com/nickfthedev/fiberHTMX/lib"
	"github.com/nickfthedev/fiberHTMX/model"
	"golang.org/x/crypto/bcrypt"
)

// PAGE
func UpdateEmailVerify(c *fiber.Ctx) error {
	key := c.Params("key")
	// Check Token
	cm := new(model.ChangeMail)
	db.DB.Where("key = ?", key).First(&cm)
	if cm.ID == 0 {
		return c.Render("user/updateemailverify", fiber.Map{"ErrorMessage": "Invalid Key."})
	}
	// Token expired
	if time.Now().After(cm.UpdatedAt.Add(60 * time.Minute)) {
		return c.Render("user/updateemailverify", fiber.Map{"ErrorMessage": "Your Key is expired.", "ErrorCode": "Please try again and verify your email within one hour."})
	}
	user := new(model.User)
	db.DB.Where("uuid = ?", cm.UserUUID).First(&user)
	if user.ID == 0 {
		return c.Render("user/updateemailverify", fiber.Map{"ErrorMessage": "User not found."})
	}
	user.Email = cm.NewMail
	db.DB.Save(&user)

	// Delete entry from ResetPassword Table
	db.DB.Unscoped().Delete(&cm)
	return c.Render("user/updateemailverify", fiber.Map{"SuccessMessage": "Email verified successfully."})
}

// HTMX
func UpdateEmail(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Email can't be empty!"}, "common/empty")
	}
	// Check if Address already exist
	u := new(model.User)
	db.DB.Where("email = ?", email).First(&u)
	if u.ID != 0 {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Email already in use!"}, "common/empty")
	}
	// Create a db entry for request
	cm := new(model.ChangeMail)
	db.DB.Where("user_uuid = ?", c.Locals("UUID")).First(&cm)
	// Create database entry
	if cm.ID == 0 { // No entry for now create a new one
		cm.UserUUID = c.Locals("UUID").(uuid.UUID)
		cm.Key = uuid.New()
		cm.NewMail = email
		db.DB.Create(&cm)
	} else { // Update the existing entry for password reset
		cm.Key = uuid.New()
		cm.NewMail = email
		db.DB.Save(&cm)
	}
	// Send mail with link
	msg := "Here is your link for changing your password for " + lib.Config.AppName + "!+<br><br> Please click the link below to verify your new email address: <br><br> <a href=\"https://" + lib.Config.Host + "/user/verifyemail/" + cm.Key.String() + "\">Verify your email address</a><br><br><i>The link is only valid for 1 hour</i>"
	errMail := lib.SendEmail(email, "Verify your new email address | "+lib.Config.AppName, msg, "text/html")
	if errMail != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to send Mail. Please try again"}, "common/empty")
	}
	return c.Render("common/success", fiber.Map{
		"SuccessMessage": "Your Email address change is now in process.",
		"SuccessCode":    "We sent you an email. Please verify your new email address",
	}, "common/empty")
}

// HTMX
func UpdateUserPassword(c *fiber.Ctx) error {
	userpw := new(model.UserChangePassword)
	if err := c.BodyParser(userpw); err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Could not parse the request body!"}, "common/empty")
	}
	if userpw.NewPassword != userpw.ConfirmPassword {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "The new password and the confirm password do not match. Please try again!"}, "common/empty")
	}
	//Check if password has at least 6 characters, one uppercase, one lowercase and one number
	password := userpw.NewPassword
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
	// Get user data from DB
	var user model.User
	if err := db.DB.Where("UUID = ?", c.Locals("UUID")).First(&user).Error; err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "User not found in DB"}, "common/empty")
	}
	// Check Password against Database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userpw.OldPassword)); err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Your current password is not correct"}, "common/empty")
	}
	//Hash Password
	hash, err := bcrypt.GenerateFromPassword([]byte(userpw.NewPassword), 10)
	if err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to hash password"}, "common/empty")
		//return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to hash password", "data": err.Error()})
	}
	user.Password = string(hash)
	// Save the updated user details
	res := db.DB.Save(&user)
	if res.Error != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to update user", "ErrorCode": res.Error.Error()}, "common/empty")

	}
	return c.Render("common/success", fiber.Map{"SuccessMessage": "Password has been changed successfully"}, "common/empty")
}

// HTMX
func UpdateUserProfile(c *fiber.Ctx) error {
	userid := c.Locals("ID")
	// Parse the form data
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Could not parse the request body!"}, "common/empty")
	}

	// Find the user in the database
	dbUser := new(model.User)
	db.DB.Where("id = ?", userid).First(&dbUser)

	// Check if the username already exist if changed
	if dbUser.Username != user.Username {
		tempUser := new(model.User)
		db.DB.Where("username = ?", user.Username).First(&tempUser)
		if tempUser.ID != 0 {
			return c.Render("common/error", fiber.Map{"ErrorMessage": "Username already exists!"}, "common/empty")
		}
	}

	// Update the user's details
	dbUser.Name = user.Name
	dbUser.Username = user.Username

	// Save the updated user details
	res := db.DB.Save(&dbUser)
	if res.Error != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to update user", "ErrorCode": res.Error.Error()}, "common/empty")

	}
	return c.Render("common/success", fiber.Map{"SuccessMessage": "User profile updated successfully"}, "common/empty")
}

// PAGE
func RenderUpdateUser(c *fiber.Ctx) error {
	// Create a new instance of the User model
	user := new(model.User)
	// Query the database for a user with the UUID provided in the request context
	db.DB.Where("uuid = ?", c.Locals("UUID")).First(&user)

	// Create a new instance of the UserSafe model
	userSafe := new(model.UserSafe)
	// Convert the User instance to JSON
	userJSON, _ := json.Marshal(user)
	//Unmarshall userJSON to userSafe instance without Password
	_ = json.Unmarshal(userJSON, &userSafe)

	//fmt.Printf("%+v\n", userSafe)
	return c.Render("user/updateprofile", fiber.Map{"user": userSafe})
}

// PAGE
func VerifyUser(c *fiber.Ctx) error {
	// Get the UUID from the URL
	uuid := c.Params("uuid")

	// Find the user with the given UUID
	user := new(model.User)
	db.DB.Where("uuid = ?", uuid).First(&user)
	fmt.Println(user.ID)
	// If the user is not found, return an error
	if user.ID == 0 {
		return c.Render("auth/login", fiber.Map{"ErrorMessage": "User not found"})
	}

	// If the user is already verified, return an error
	if user.Verified {
		return c.Render("auth/login", fiber.Map{"ErrorMessage": "User already verified"})
	}

	// Otherwise, verify the user and save to the database
	user.Verified = true
	db.DB.Save(&user)

	// Return a success message
	return c.Render("auth/login", fiber.Map{"SuccessMessage": "User verified successfully"})
}

// HTMX
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
	user.UUID = uuid.New()
	result := db.DB.Create(&user) // pass pointer of data to Create
	if result.Error != nil {      //Error handling
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
			return c.Render("common/error", fiber.Map{"ErrorMessage": "Registration failed", "ErrorCode": "E-Mail is already in use "}, "common/empty")
		}
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to create user", "ErrorCode": result.Error.Error()}, "common/empty")
	}
	// Send Verification Mail
	msg := "Welcome to " + lib.Config.AppName + "!<br><br> Please verify your account and click the link: <a href=\"https://" + lib.Config.Host + "/user/verify/" + user.UUID.String() + "\">Verify Account</a>"
	errMail := lib.SendEmail(user.Email, "Activate your Account | "+lib.Config.AppName, msg, "text/html")
	if errMail != nil {
		return c.Render("common/error", fiber.Map{"ErrorMessage": "Failed to send activation Mail, but registered successfully."}, "common/empty")
	}
	//If no error send back ok
	return c.Render("common/success", fiber.Map{"SuccessMessage": "Registered successfully", "SuccessCode": "We send you an email to verify your account. After verification you can login"}, "common/empty")
}

// INIT APP
func CreateStandardAdminUser() {
	//Map input to user Model
	user := new(model.User)
	user.Username = "admin"
	user.Name = "admin"
	user.Email = "admin@admin.com"
	user.Password = "password"
	user.Verified = true
	user.UUID = uuid.New()
	//Hash Password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)
	db.DB.Create(&user) // pass pointer of data to Create
	fmt.Println("\n=============================\n", "No User in database. Created user 'admin@admin.com' with password 'password'!", "\n=============================")
}
