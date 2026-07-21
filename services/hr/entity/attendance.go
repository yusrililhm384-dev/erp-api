package entity

import (
	"time"

	"enterprise_resource_planning/internal/util"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	Date          util.JSONDate `json:"date" gorm:"type:date;not null;uniqueIndex:idx_employee_date"`
	CheckIn       *time.Time    `json:"check_in" gorm:"type:timestamptz;index"`
	CheckOut      *time.Time    `json:"check_out" gorm:"type:timestamptz;index"`
	Status        string        `json:"status" gorm:"type:varchar(25);not null;check:,status <> 'present','absent','late','sick'"`
	Note          *string       `json:"note" gorm:"type:text"`
	AttachmentKey *string       `json:"attachment_key" gorm:"type:text"`
	EmployeeID    uint          `json:"employee_id" gorm:"uniqueIndex:idx_employee_date"`
	Employee      *Employee     `json:"employee" gorm:"foreignKey:EmployeeID"`
}
