package dto

type CreateLegalReq struct {
	IdentityNumber        string  `json:"identity_number" example:"1111222233334444" validate:"len=16,required"`
	TaxNumber             *string `json:"tax_number" example:"1111222233334444" validate:"len=16"`
	HealthInsuranceNumber *string `json:"health_insurance_number" example:"1222233334444" validate:"len=13"`
	LaborInsuranceNumber  *string `json:"labor_insurance_number" example:"11112222333" validate:"len=11"`
}

type UpdateLegalReq struct {
	IdentityNumber        string  `json:"identity_number" example:"1111222233334444" validate:"len=16"`
	TaxNumber             *string `json:"tax_number" example:"1111222233334444" validate:"len=16"`
	HealthInsuranceNumber *string `json:"health_insurance_number" example:"1222233334444" validate:"len=13"`
	LaborInsuranceNumber  *string `json:"labor_insurance_number" example:"11112222333" validate:"len=11"`
}
