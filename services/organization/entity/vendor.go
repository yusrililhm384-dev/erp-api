package entity

import (
	"gorm.io/gorm"
)

type Vendor struct {
	gorm.Model

	VendorTypeID uint        `json:"vendor_type_id" gorm:"index;not null"`
	VendorType   *VendorType `json:"vendor_type" gorm:"foreignKey:VendorTypeID"`

	CountryID uint     `json:"country_id" gorm:"index;not null"`
	Country   *Country `json:"country" gorm:"foreignKey:CountryID"`

	LanguageID uint      `json:"language_id" gorm:"index;not null"`
	Language   *Language `json:"language" gorm:"foreignKey:LanguageID"`

	CurrencyID uint      `json:"currency_id" gorm:"index;not null"`
	Currency   *Currency `json:"currency" gorm:"foreignKey:CurrencyID"`

	DefaultPaymentTermID uint         `json:"default_payment_term_id" gorm:"index;not null"`
	DefaultPaymentTerm   *PaymentTerm `json:"default_payment_term" gorm:"foreignKey:DefaultPaymentTermID"`

	Code        string  `json:"code" gorm:"type:varchar(20);not null;uniqueIndex"`
	Name        string  `json:"name" gorm:"type:varchar(100);not null"`
	Description string  `json:"description" gorm:"type:text;not null"`
	TaxNumber   *string `json:"tax_number" gorm:"type:varchar(30);uniqueIndex"`
	Address     string  `json:"address" gorm:"type:text;not null"`
	Location    string  `json:"location" gorm:"type:geography(POINT,4326);not null"`

	Contacts []*VendorContact `json:"contacts" gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;"`
}

type VendorContact struct {
	gorm.Model
	Name      string  `json:"name" gorm:"type:varchar(50)"`
	Email     string  `json:"email" gorm:"type:varchar(100);not null;uniqueIndex:idx_vendor_email_active,where:deleted_at IS NULL"`
	Phone     string  `json:"phone" gorm:"type:varchar(15);not null;uniqueIndex:idx_vendor_phone_active,where:deleted_at IS NULL"`
	Position  string  `json:"position" gorm:"type:varchar(100);not null"`
	VendorID  uint    `json:"vendor_id" gorm:"not null;uniqueIndex:idx_vendor_email_active;uniqueIndex:idx_vendor_phone_active"`
	Vendor    *Vendor `json:"vendor" gorm:"foreignKey:VendorID"`
	IsPrimary bool    `json:"is_primary" gorm:"type:boolean;not null"`
}

type VendorType struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
}

func (v *Vendor) BeforeDelete(tx *gorm.DB) error {
	return tx.Where("vendor_id = ?").Delete(&VendorContact{}).Error
}
