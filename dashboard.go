package main

import (
	"test-helloworld-ci-cd/database"

	"github.com/gofiber/fiber/v2"
)

func ShowDashboard(c *fiber.Ctx) error {
	var reservations []database.Reservation
	database.DB.Find(&reservations)
	return c.Render("dashboard", fiber.Map{
		"Reservations": reservations,
	}, "layouts/main") // of gewoon .SendFile() als je geen templates gebruikt
}

func DeleteReservation(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&database.Reservation{}, id)
	return c.Redirect("/dashboard")
}

func UpdateReservation(c *fiber.Ctx) error {
	id := c.Params("id")
	var updateData database.Reservation

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).SendString("Fout bij verwerken")
	}

	var res database.Reservation
	database.DB.First(&res, id)

	res.Name = updateData.Name
	res.Date = updateData.Date
	res.Time = updateData.Time
	res.Status = updateData.Status

	database.DB.Save(&res)
	return c.Redirect("/dashboard")
}
