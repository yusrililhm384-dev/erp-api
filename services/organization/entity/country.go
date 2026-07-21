package entity

import "gorm.io/gorm"

type Country struct {
	gorm.Model
	Code string `json:"code" gorm:"uniqueIndex;type:varchar(2);not null"`
	Name string `json:"name" gorm:"type:varchar(50);not null"`
}
