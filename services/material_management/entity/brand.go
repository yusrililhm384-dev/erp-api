package entity

import "gorm.io/gorm"

type Brand struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(50);not null"`
}
