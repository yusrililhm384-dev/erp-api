package entity

import (
	"time"

	"enterprise_resource_planning/services/organization/entity"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name      string              `json:"name" gorm:"type:varchar(100);not null"`
	Email     string              `json:"email" gorm:"uniqueIndex:idx_active_email,where:deleted_at is null;type:varchar(100);not null"`
	Phone     string              `json:"phone" gorm:"uniqueIndex:idx_active_phone,where:deleted_at is null;type:varchar(15);not null"`
	Gender    bool                `json:"gender" gorm:"not null;type:boolean"`
	Address   string              `json:"address" gorm:"type:text;not null"`
	BirthDate time.Time           `json:"birth_date" gorm:"type:date;not null"`
	JoinDate  time.Time           `json:"join_date" gorm:"type:date;not null"`
	Legal     *EmployeeLegal      `json:"legal" gorm:"constraint:OnDelete:CASCADE;OnUpdate:CASCADE"`
	User      *User               `json:"user" gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE"`
	Mutations []*EmployeeMutation `json:"mutations" gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE"`
}

type EmployeeLegal struct {
	gorm.Model
	EmployeeID            uint    `json:"employee_id" gorm:"uniqueIndex:idx_active_employee,where:deleted_at is null;not null"`
	IdentityNumber        string  `json:"identity_number" gorm:"uniqueIndex:idx_active_identity,where:deleted_at is null;type:varchar(16);not null"`
	TaxNumber             *string `json:"tax_number" gorm:"uniqueIndex:idx_active_tax,where:deleted_at is null;type:varchar(16)"`
	HealthInsuranceNumber *string `json:"health_insurance_number" gorm:"uniqueIndex:idx_active_health_insurance,where:deleted_at is null;type:varchar(13)"`
	LaborInsuranceNumber  *string `json:"labor_insurance_number" gorm:"uniqueIndex:idx_active_labor_insurance,where:deleted_at is null;type:varchar(11)"`
}

type EmployeeMutation struct {
	gorm.Model
	EmployeeID   uint               `json:"employee_id" gorm:"index;not null"`
	BranchID     uint               `json:"branch_id" gorm:"index;not null"`
	PositionID   uint               `json:"position_id" gorm:"index;not null"`
	TypeID       uint               `json:"type_id" gorm:"index;notnull"`
	DepartmentID uint               `json:"department_id" gorm:"index;not null"`
	StartDate    time.Time          `json:"start_date" gorm:"type:date;not null"`
	EndDate      *time.Time         `json:"end_date,omitempty" gorm:"type:date;index"`
	Employee     *Employee          `json:"employee" gorm:"foreignKey:EmployeeID"`
	Branch       *entity.Branch     `json:"branch" gorm:"foreignKey:BranchID"`
	Position     *entity.Position   `json:"position" gorm:"foreignKey:PositionID"`
	Type         *EmployeeType      `json:"type" gorm:"foreignKey:TypeID"`
	Department   *entity.Department `json:"department" gorm:"foreignKey:DepartmentID"`
}

type EmployeeType struct {
	gorm.Model
	Name        string              `json:"name" gorm:"type:varchar(50);not null"`
	Description string              `json:"description" gorm:"type:text;not null"`
	Mutations   []*EmployeeMutation `json:"mutations" gorm:"constraint:OnDelete:RESTRICT,OnUpdate:CASCADE;"`
}

func (e *Employee) BeforeDelete(tx *gorm.DB) error {
	if err := tx.Where("employee_id = ?", e.ID).Delete(&EmployeeLegal{}).Error; err != nil {
		return err
	}

	if err := tx.Where("employee_id = ?", e.ID).Delete(&EmployeeMutation{}).Error; err != nil {
		return err
	}

	return tx.Where("employee_id = ?", e.ID).Delete(&User{}).Error
}
