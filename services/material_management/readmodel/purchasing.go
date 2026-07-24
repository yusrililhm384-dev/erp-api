package readmodel

import "time"

type PurchasingOrganizationListResponse struct {
	Meta *MetaPagination               `json:"meta"`
	List []*PurchasingOrganizationList `json:"list"`
}

type PurchasingOrganizationList struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Company   string    `json:"company"`
	Country   string    `json:"country"`
	Language  string    `json:"language"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PurchasingOrganizationDetail struct {
	Id          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Company     string    `json:"company"`
	Country     string    `json:"country"`
	Language    string    `json:"language"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PurchasingGroupListResponse struct {
	Meta *MetaPagination        `json:"meta"`
	List []*PurchasingGroupList `json:"list"`
}

type PurchasingGroupList struct {
	Id                     uint      `json:"id"`
	Code                   string    `json:"code"`
	Name                   string    `json:"name"`
	PurchasingOrganization string    `json:"purchasing_organization"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type PurchasingGroupDetail struct {
	Id                     uint                     `json:"id"`
	Code                   string                   `json:"code"`
	Name                   string                   `json:"name"`
	Description            string                   `json:"description"`
	PurchasingOrganization string                   `json:"purchasing_organization"`
	PurchasingGroupMembers []*PurchasingGroupMember `json:"list"`
	CreatedAt              time.Time                `json:"created_at"`
	UpdatedAt              time.Time                `json:"updated_at"`
}

type PurchasingGroupMember struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Position string `json:"position"`
}
