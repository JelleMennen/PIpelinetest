package auth

import (
	"github.com/gofiber/fiber/v2"
)

var isLoggedIn bool = false

func ShowLoginPage(c *fiber.Ctx) error {
	return c.SendFile("views/login.html")
}

func HandleLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "admin" && password == "test123" {
		isLoggedIn = true
		return c.Redirect("/dashboard")
	}
	return c.SendString("<p style='color: red;'>Ongeldige inloggegevens</p>")
}

func HandleLogout(c *fiber.Ctx) error {
	isLoggedIn = false
	return c.Redirect("/")
}

func RequireLogin(c *fiber.Ctx) error {
	if !isLoggedIn {
		return c.Redirect("/login")
	}
	return c.Next()
}
