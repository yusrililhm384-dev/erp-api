package entity

import "gorm.io/gorm"

type MaterialGroup struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null;type:varchar(100)"`
	Description string `json:"description" gorm:"not null;type:text"`
}
