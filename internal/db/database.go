package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo-planning/internal/config"
)

type Config struct {
	Type     string `yaml:"type"` // sqlite or postgres
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func NewConnection() (*gorm.DB, error) {
	config, err := config.Load()
	if err != nil {
		return nil, err
	}

	switch config.Database.Driver {
	case "sqlite":
		return gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", config.Database.Name)), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Name, config.Database.SSLMode)
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Database.Driver)
	}

}
