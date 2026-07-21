package entity

import "gorm.io/gorm"

type Warehouse struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(50);not null"`
	Description string `json:"desciption" gorm:"type:text;not null"`
	Address     string `json:"address" gorm:"type:text;not null"`
	Location    string `json:"location" gorm:"type:geography(POINT,4326);not null"`
}
