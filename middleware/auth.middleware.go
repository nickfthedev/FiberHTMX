package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/nickfthedev/fiberHTMX/db"
	"github.com/nickfthedev/fiberHTMX/lib"
	"github.com/nickfthedev/fiberHTMX/model"
)

type Authorization struct {
	Token string `json:"token"`
}

func CheckLogin(c *fiber.Ctx) (model.User, error) {
	var user model.User
	// Check if tokenstring is in cookie
	tokenString := c.Cookies("Authorization", "")
	//Check if Token was sent in request
	if tokenString == "" {
		t := new(Authorization)
		if err := c.BodyParser(t); err != nil {
			return user, fiber.ErrForbidden
		}
	}
	if tokenString == "" {
		return user, fiber.ErrForbidden
	}
	// Validate
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")

		return []byte(lib.Config.TokenSecret), nil
	})
	if err != nil {
		log.Println(err.Error())
		return user, fiber.ErrBadGateway
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return user, fiber.ErrForbidden
		}
		//Find User
		db.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			return user, fiber.ErrForbidden
		}
		return user, nil

	} else {
		return user, fiber.ErrForbidden
	}
}

// Sets c.Locals for Name, ID and IsLoggedIn for use in templates
func IsLoggedIn(c *fiber.Ctx) error {
	user, _ := CheckLogin(c)
	if user.ID != 0 {
		c.Locals("Name", user.Name)
		c.Locals("UUID", user.UUID)
		c.Locals("ID", user.ID)
		c.Locals("IsLoggedIn", true)
	}
	return c.Next()
}

// Middleware Access for Guest only for example register and login
func GuestOnly(c *fiber.Ctx) error {
	user, _ := CheckLogin(c)
	if user.ID == 0 {
		return c.Next()
	}
	return c.Redirect("/")
}

// Access for logged in users only
func LoginRequired(c *fiber.Ctx) error {
	_, err := CheckLogin(c)
	if err != nil {
		return c.Redirect("/")
	}
	return c.Next()
}
