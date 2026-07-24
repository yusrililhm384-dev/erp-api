package dto

type CreatePurchasingOrganizationReq struct {
	Code         string `json:"code" example:"DOME" validate:"required"`
	Name         string `json:"name" example:"Domestic" validate:"required"`
	Description  string `json:"description" example:"Domestic Purchasing" validate:"required"`
	CompanyId    uint   `json:"company_id" example:"1" validate:"required"`
	CountryId    uint   `json:"country_id" example:"1" validate:"required"`
	LanguageId   uint   `json:"language_id" example:"1" validate:"required"`
	CurrencyIdId uint   `json:"currency_id" example:"1" validate:"required"`
}

type UpdatePurchasingOrganizationReq struct {
	Code         string `json:"code" example:"DOME"`
	Name         string `json:"name" example:"Domestic"`
	Description  string `json:"description" example:"Domestic Purchasing"`
	CompanyId    uint   `json:"company_id" example:"1"`
	CountryId    uint   `json:"country_id" example:"1"`
	LanguageId   uint   `json:"language_id" example:"1"`
	CurrencyIdId uint   `json:"currency_id" example:"1"`
}

type CreatePurchasingGroupReq struct {
	Code                     string                      `json:"code" example:"1000" validate:"required"`
	Name                     string                      `json:"name" example:"IT" validate:"required"`
	Description              string                      `json:"description" example:"IT Buyer" validate:"required"`
	PurchasingOrganizationId uint                        `json:"puchasing_organization_Id" example:"1" validate:"required"`
	Members                  []*PurchasingGroupMemberReq `json:"purchasing_group_members,omitempty"`
}

type UpdatePurchasingGroupReq struct {
	Code                     string `json:"code" example:"1000"`
	Name                     string `json:"name" example:"IT"`
	Description              string `json:"description" example:"IT Buyer"`
	PurchasingOrganizationId uint   `json:"purchasing_organization_id" example:"1"`
}

type PurchasingGroupMemberReq struct {
	PurchasingGroupId uint `json:"purchasing_group_id" example:"1"`
	EmployeeId        uint `json:"employee_id" example:"1"`
}
