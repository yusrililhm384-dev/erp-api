package entity

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model

	Code        string  `json:"code" gorm:"not null;uniqueIndex;type:varchar(100)"`
	Name        string  `json:"name" gorm:"not null;type:varchar(100)"`
	Description string  `json:"description" gorm:"not null;type:text"`
	Phone       string  `json:"phone" gorm:"not null;uniqueIndex;type:varchar(15)"`
	Email       string  `json:"email" gorm:"not null;uniqueIndex;type:varchar(100)"`
	Website     *string `json:"website" gorm:"varchar(100)"`
	Address     string  `json:"address" gorm:"not null;type:text"`
	Location    string  `json:"location" gorm:"type:geography(POINT,4326);not null"`

	CountryID uint     `json:"country_id" gorm:"not null;index"`
	Country   *Country `json:"country" gorm:"foreignKey:CountryID"`

	CurrencyID uint      `json:"currency_id" gorm:"not null;index"`
	Currency   *Currency `json:"currency" gorm:"foreignKey:CurrencyID"`

	LanguageID uint      `json:"language_id" gorm:"not null;index"`
	Language   *Language `json:"language" gorm:"foreignKey:LanguageID"`

	Branches []*Branch `json:"branches" gorm:"foreignKey:CompanyID"`
}
