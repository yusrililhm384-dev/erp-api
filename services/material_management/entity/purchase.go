package entity

import (
	"enterprise_resource_planning/services/hr/entity"

	"gorm.io/gorm"
)

type PurchaseGroup struct {
	gorm.Model
	Code        string `json:"code" grom:"not null;type:varchar(100);uniqueIndex"`
	Name        string `json:"name" grom:"not null;type:varchar(100)"`
	Description string `json:"description" gorm:"not null;type:text"`
}

type PurchaseGroupRole struct {
	gorm.Model
	Name string `json:"name" grom:"not null;type:varchar(100)"`
}

type PurchaseGroupMember struct {
	gorm.Model

	PurchaseGroupID uint           `json:"purchase_group_id" gorm:"not null;uniqueIndex:uq_idx_purchase_group_members_purchase_employee_role_id"`
	PurchaseGroup   *PurchaseGroup `json:"purchase_group" gorm:"foreignKey:PurhcaseGroupID"`

	EmployeeID uint             `json:"employee_Id" gorm:"not null;uniqueIndex:uq_idx_purchase_group_members_purchase_employee_role_id"`
	Employee   *entity.Employee `json:"employee" gorm:"foreignKey:EmployeeID"`

	RoleID uint               `json:"role_id" gorm:"not null;uniqueIndex:uq_idx_purchase_group_members_purchase_employee_role_id"`
	Role   *PurchaseGroupRole `json:"role" gorm:"foreignKey:RoleID"`
}

type PurchaseOrganization struct {
	gorm.Model
}
