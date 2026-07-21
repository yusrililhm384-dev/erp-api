package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormPg() (*gorm.DB, error) {
	dsn := ""

	pg, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return pg, nil
}
