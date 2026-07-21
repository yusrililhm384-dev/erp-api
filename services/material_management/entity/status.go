package entity

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Code string `json:"code" gorm:"uniqueIndex;not null;type:varchar(100)"`
	Name string `json:"name" gorm:"type:varchar(100);not null"`
}
