package entity

import (
	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(50);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
}

func (d *Department) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.Model(&Position{}).Where("department_id = ?", d.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errs.ErrDepartmentHasActivePosition
	}

	return nil
}
