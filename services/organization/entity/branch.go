package entity

import (
	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(50);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	Address     string `json:"address" gorm:"type:text;not null"`
	Location    string `json:"location" gorm:"type:geography(POINT,4326);not null"`

	CompanyID uint     `json:"company_id" gorm:"index;not null"`
	Company   *Company `json:"company" gorm:"foreignKey:CompanyID"`
}

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (b *Branch) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.Table("employee_mutations").Where("branch_id = ?", b.ID).Where("deleted_at is null").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errs.ErrBranchHasActiveEmployees
	}

	return nil
}
