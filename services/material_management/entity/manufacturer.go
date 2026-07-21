package entity

import "gorm.io/gorm"

type Manufacturer struct {
	gorm.Model
	Code string `json:"code" gorm:"uniqueIndex;not null;type:varchar(100)"`
	Name string `json:"name" gorm:"not null;type:varchar(100)"`
}
