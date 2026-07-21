package entity

import (
	"enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

type MaterialUnit struct {
	gorm.Model
	Code                   string                `json:"code" gorm:"not null;type:varchar(15);uniqueIndex"`
	Name                   string                `json:"name" gorm:"not null;type:varchar(100)"`
	Description            string                `json:"description" gorm:"not null;type:text"`
	MaterialUnitCategoryID uint                  `json:"material_unit_category_id"`
	MaterialUnitCategory   *MaterialUnitCategory `json:"material_unit_category" gorm:"foreignKey:MaterialUnitCategoryID"`
}

type MaterialUnitCategory struct {
	gorm.Model
	Code          string          `json:"code" gorm:"not null;type:varchar(15);uniqueIndex"`
	Name          string          `json:"name" gorm:"not null;type:varchar(100)"`
	MaterialUnits []*MaterialUnit `json:"material_units" gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
}

func (muc *MaterialUnitCategory) BeforeDelete(tx *gorm.DB) error {
	var count int64

	if err := tx.WithContext(tx.Statement.Context).Model(&MaterialUnit{}).Where("material_unit_category_id = ?", muc.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.ErrMaterialUnitCategoryHasMaterialUnit
	}

	return nil
}
