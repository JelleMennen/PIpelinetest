package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Globale database variabele
var DB *gorm.DB

// Connectie maken met de MySQL database
func ConnectDatabase() {
	dsn := "reservering:WelkomWelkom040!@tcp(135.225.96.153:3306)/reserveringen_db?parseTime=true"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Verbonden met de MySQL database!")

	DB = database
}

// Database model van reserveringen
type Reservation struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Status string `json:"status"`
}

// Automatische migratie van de database-tabellen
func MigrateDatabase() {
	DB.AutoMigrate(&Reservation{})
	log.Println("Database tabellen gemigreerd!")
}
