package entity

import (
	"enterprise_resource_planning/services/hr/entity"

	"gorm.io/gorm"
)

type PurchasingGroup struct {
	gorm.Model
	Code                     string `json:"code" grom:"not null;type:varchar(100);uniqueIndex"`
	Name                     string `json:"name" grom:"not null;type:varchar(100)"`
	Description              string `json:"description" gorm:"not null;type:text"`

	PurchasingOrganizationID uint   `json:"purchasing_organization_id" gorm:"index;not null"`
	PurchasingOrganization *PurchasingOrganization `json:"purchasing_organization"`
}

type PurchasingGroupRole struct {
	gorm.Model
	Name string `json:"name" grom:"not null;type:varchar(100)"`
}

type PurchasingGroupMember struct {
	gorm.Model

	PurchasingGroupID uint             `json:"purchasing_group_id" gorm:"not null;uniqueIndex:uq_idx_purchasing_group_members_purchasing_employee_role_id"`
	PurchasingGroup   *PurchasingGroup `json:"purchasing_group" gorm:"foreignKey:PurhcaseGroupID"`

	EmployeeID uint             `json:"employee_Id" gorm:"not null;uniqueIndex:uq_idx_purchasing_group_members_purchasing_employee_role_id"`
	Employee   *entity.Employee `json:"employee" gorm:"foreignKey:EmployeeID"`

	RoleID uint                 `json:"role_id" gorm:"not null;uniqueIndex:uq_idx_purchasing_group_members_purchasing_employee_role_id"`
	Role   *PurchasingGroupRole `json:"role" gorm:"foreignKey:RoleID"`
}

type PurchasingOrganization struct {
	gorm.Model
	Code        string `json:"code" grom:"not null;type:varchar(100);uniqueIndex"`
	Name        string `json:"name" grom:"not null;type:varchar(100)"`
	Description string `json:"description" gorm:"not null;type:text"`
}
