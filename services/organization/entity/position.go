package entity

import (
	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

type PositionCategory struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(50);not null"`
}

type Position struct {
	gorm.Model
	Name               string            `json:"name" gorm:"type:varchar(100);not null"`
	Description        string            `json:"description" gorm:"type:text;not null"`
	ApprovalID         *uint             `json:"approval_id" gorm:"type:int;check:approval_id_not_equal, id <> approval_id"`
	Approval           *Position         `json:"approval" gorm:"foreignKey:ApprovalID"`
	PositionCategoryID uint              `json:"position_category_id" gorm:"index;not null"`
	PositionCategory   *PositionCategory `json:"position_category" gorm:"foreignKey:PositionCategoryID"`
}

func (pc *PositionCategory) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.Model(&Position{}).Where("position_category_id = ?", pc.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errs.ErrCategoryHasPosition
	}

	return nil
}

func (p *Position) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.Model(&Position{}).Where("approval_id = ?", p.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errs.ErrPositionHasActiveEmployees
	}

	if err := tx.Table("employee_mutations").Where("position_id = ?", p.ID).Where("deleted_at is null").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errs.ErrPositionHasActiveEmployees
	}

	return nil
}
