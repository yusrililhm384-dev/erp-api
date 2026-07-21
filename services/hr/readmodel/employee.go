package readmodel

import (
	"time"

	"enterprise_resource_planning/internal/util"
)

type EmployeeListResponse struct {
	Meta *MetaPagination `json:"meta"`
	List []*EmployeeList `json:"list"`
}

type EmployeeList struct {
	Id        uint          `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Phone     string        `json:"phone"`
	Gender    string        `json:"gender"`
	JoinDate  util.JSONDate `json:"join_date"`
	Position  string        `json:"position"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type EmployeeDetail struct {
	Id        uint          `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Phone     string        `json:"phone"`
	Gender    string        `json:"gender"`
	Address   string        `json:"address"`
	JoinDate  util.JSONDate `json:"join_date"`
	Legal     *Legal        `json:"legal"`
	Histories []*History    `json:"histories"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Legal struct {
	IdentityNumber        string  `json:"identity_number"`
	TaxNumber             *string `json:"tax_number,omitempty"`
	HealthInsuranceNumber *string `json:"health_insurance+number,omitempty"`
	LaborInsuranceNumber  *string `json:"labor_insurance_number,omitempty"`
}

type History struct {
	Position   string            `json:"position"`
	Branch     string            `json:"branch"`
	Department string            `json:"department"`
	Type       string            `json:"type"`
	StartDate  util.JSONDate     `json:"start_date"`
	EndDate    util.JSONNullDate `json:"end_date"`
}
