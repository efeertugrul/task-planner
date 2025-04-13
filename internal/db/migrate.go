package db

import (
	"gorm.io/gorm"

	"todo-planning/internal/model"
)

// AutoMigrate performs automatic schema migration for all models
func AutoMigrate(db *gorm.DB) error {
	// Enable foreign key constraints for SQLite if using SQLite database
	if db.Dialector.Name() == "sqlite" {
		db.Exec("PRAGMA foreign_keys = ON")
	}

	// Migrate the schema
	return db.AutoMigrate(
		&model.Task{},
		&model.Developer{},
		&model.Assignment{},
	)
}
