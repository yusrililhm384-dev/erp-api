package entity

import "gorm.io/gorm"

type IndustrySector struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null;type:varchar(50)"`
	Description string `json:"description" gorm:"not null;type:text"`
}
