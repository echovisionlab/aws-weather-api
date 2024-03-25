package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(config *Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(config.ConnStr()))
}
