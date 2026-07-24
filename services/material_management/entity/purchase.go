package entity

import (
	"enterprise_resource_planning/services/hr/entity"

	org "enterprise_resource_planning/services/organization/entity"

	"gorm.io/gorm"
)

type PurchasingGroup struct {
	gorm.Model
	Code        string `json:"code" gorm:"not null;type:varchar(20);uniqueIndex"`
	Name        string `json:"name" gorm:"not null;type:varchar(100)"`
	Description string `json:"description" gorm:"not null;type:text"`

	PurchasingOrganizationID uint                    `json:"purchasing_organization_id" gorm:"index;not null"`
	PurchasingOrganization   *PurchasingOrganization `json:"purchasing_organization" gorm:"foreignKey:PurchasingOrganizationID"`

	PurchasingGroupMembers []*PurchasingGroupMember `json:"purchasing_group_members" gorm:"foreignKey:PurchasingGroupID"`
}

type PurchasingGroupMember struct {
	gorm.Model

	PurchasingGroupID uint             `json:"purchasing_group_id" gorm:"not null;uniqueIndex:uq_idx_purchasing_group_members_purchasing_employee_role_id"`
	PurchasingGroup   *PurchasingGroup `json:"purchasing_group" gorm:"foreignKey:PurchasingGroupID"`

	EmployeeID uint             `json:"employee_id" gorm:"not null;uniqueIndex:uq_idx_purchasing_group_members_purchasing_employee_role_id"`
	Employee   *entity.Employee `json:"employee" gorm:"foreignKey:EmployeeID"`
}

type PurchasingOrganization struct {
	gorm.Model

	Code        string `json:"code" gorm:"not null;type:varchar(20);uniqueIndex"`
	Name        string `json:"name" gorm:"not null;type:varchar(100)"`
	Description string `json:"description" gorm:"not null;type:text"`

	CompanyID uint         `json:"company_id" gorm:"index;not null"`
	Company   *org.Company `json:"company" gorm:"foreignKey:CompanyID"`

	CountryID uint         `json:"country_id" gorm:"index;not null"`
	Country   *org.Country `json:"country" gorm:"foreignKey:CountryID"`

	LanguageID uint          `json:"language_id" gorm:"index;not null"`
	Language   *org.Language `json:"language" gorm:"foreignKey:LanguageID"`

	CurrencyID uint          `json:"currency_id" gorm:"index;not null"`
	Currency   *org.Currency `json:"currency" gorm:"foreignKey:CurrencyID"`

	PurchasingGroups []*PurchasingGroup `json:"purchasing_groups" gorm:"foreignKey:PurchasingOrganizationID"`
}
