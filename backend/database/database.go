package database

import (
	"log"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("flight_data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Airport{}, &models.Aircraft{}, &models.Flight{})
	if err != nil {
		log.Fatal("Failed to migrate database")
		return nil, err
	}

	return db, nil
}
