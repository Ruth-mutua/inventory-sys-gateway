package main

import (
	"log"

	"go-inventory-system/shared"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// initDatabase initializes the database connection and runs migrations
func initDatabase(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the Order model
	if err := db.AutoMigrate(&shared.Order{}); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
