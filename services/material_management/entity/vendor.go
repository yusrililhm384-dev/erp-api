package entity

import "gorm.io/gorm"

type Vendor struct {
	gorm.Model

	VendorTypeID uint        `json:"vendor_type_id" gorm:"index;not null"`
	VendorType   *VendorType `json:"vendor_type" gorm:"foreignKey:VendorTypeID"`

	Code        string           `json:"code" gorm:"type:varchar(20);not null;uniqueIndex"`
	Name        string           `json:"name" gorm:"type:varchar(100);not null"`
	Description string           `json:"description" gorm:"type:text;not null"`
	Address     string           `json:"address" gorm:"type:text;not null"`
	Location    string           `json:"location" gorm:"type:geography(POINT,4326);not null"`
	Contacts    []*VendorContact `json:"contacts" gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;"`
}

type VendorContact struct {
	gorm.Model
	Name     string  `json:"name" gorm:"type:varchar(50)"`
	Email    string  `json:"email" gorm:"type:varchar(100);not null;uniqueIndex:idx_vendor_email_active,where:deleted_at IS NULL"`
	Phone    string  `json:"phone" gorm:"type:varchar(15);not null;uniqueIndex:idx_vendor_phone_active,where:deleted_at IS NULL"`
	VendorID uint    `json:"vendor_id" gorm:"not null;uniqueIndex:idx_vendor_email_active;uniqueIndex:idx_vendor_phone_active"`
	Vendor   *Vendor `json:"vendor" gorm:"foreignKey:VendorID"`
}

type VendorType struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
}

func (v *Vendor) BeforeDelete(tx *gorm.DB) error {
	return tx.Where("vendor_id = ?").Delete(&VendorContact{}).Error
}
