package entity

import "gorm.io/gorm"

type Division struct {
	gorm.Model
	Name string `json:"name" gorm:"not null;type:varchar(100)"`
}
