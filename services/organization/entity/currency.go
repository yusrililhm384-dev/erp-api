package entity

import "gorm.io/gorm"

type Currency struct {
	gorm.Model

	Code string `json:"code" gorm:"type:varchar(5);not null"`
	Name string `json:"name" gorm:"type:varchar(100);not null"`

	CountryID uint       `json:"country_id" gorm:"index;not null"`
	Country   []*Country `json:"country" gorm:"foreignKey:CountryID"`
}
