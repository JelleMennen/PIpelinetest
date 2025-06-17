package main

import (
	"strconv"
	"test-helloworld-ci-cd/database"

	"test-helloworld-ci-cd/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	database.ConnectDatabase()

	// Snel testen of je kan migreren en schrijven
	//
	err := database.DB.AutoMigrate(&database.Reservation{})
	if err != nil {
		panic("Mislukt om tabel aan te maken: " + err.Error())
	}
	// Voeg testdata toe
	test := database.Reservation{
		Date:   "2025-06-16",
		Time:   "12:00",
		Status: "Bevestigd",
		Name:   "TestAzure",
	}
	if err := database.DB.Create(&test).Error; err != nil {
		panic("Testdata toevoegen mislukt: " + err.Error())
	}
	// Je ziet direct of het werkt: bij een error stopt de app meteen.

	database.MigrateDatabase()

	database.ConnectDatabase()
	database.MigrateDatabase()

	app := fiber.New()
	app.Use(logger.New()) //middleware voor logging

	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("views/index.html")
	})

	// Serve de indexpagina
	app.Get("/reservering", func(c *fiber.Ctx) error {
		return c.SendFile("views/reservering.html")
	})

	// API: Haal alle reserveringen op
	app.Get("/reservations", func(c *fiber.Ctx) error {
		var reservations []database.Reservation
		database.DB.Find(&reservations)
		return c.JSON(reservations)
	})

	// API: Maak een nieuwe reservering aan
	app.Post("/reservations", func(c *fiber.Ctx) error {
		var reservation database.Reservation

		if err := c.BodyParser(&reservation); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Ongeldige invoer"})
		}

		database.DB.Create(&reservation)

		accept := c.Get("Accept")
		if accept == "application/json" {
			return c.Status(201).JSON(reservation)
		}
		return c.SendString("<p style='color: green;'>Reservering bevestigd voor " + reservation.Date + " om " + reservation.Time + ".</p>")
	})

	app.Get("/login", auth.ShowLoginPage)
	app.Post("/login", auth.HandleLogin)
	app.Get("/logout", auth.HandleLogout)

	app.Use("/dashboard", auth.RequireLogin)
	app.Use("/delete", auth.RequireLogin)
	app.Use("/update", auth.RequireLogin)

	//
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.SendFile("views/dashboard.html")
	})

	//Dashboard laadt html zien
	app.Get("/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		database.DB.Delete(&database.Reservation{}, id)
		return c.Redirect("/dashboard")
	})

	// Reservering verwijderen
	app.Get("/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		database.DB.Delete(&database.Reservation{}, id)
		return c.Redirect("/dashboard")
	})

	//Reserering updaten
	app.Post("/update/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).SendString("Ongeldige ID")
		}

		var updated database.Reservation
		if err := c.BodyParser(&updated); err != nil {
			return c.Status(400).SendString("Fout bij verwerken formulier")
		}

		var original database.Reservation
		database.DB.First(&original, id)
		original.Status = updated.Status

		database.DB.Save(&original)
		return c.Redirect("/dashboard")
	})

	app.Listen(":80")
}
