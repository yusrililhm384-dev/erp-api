package entity

import "gorm.io/gorm"

type Language struct {
	gorm.Model
	Code string `json:"code" gorm:"uniqueIndex;type:varchar(100);not null"`
	Name string `json:"name" gorm:"type:varchar(100);not null"`
}
