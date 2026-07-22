package dto

type CreatePaymentTermReq struct {
	Code        string `json:"code" example:"NET30" validate:"required"`
	Name        string `json:"name" example:"Net 30" validate:"required"`
	Description string `json:"description" example:"This means the buyer has 30 days from the invoice date to pay" validate:"required"`
}

type UpdatePaymentTermReq struct {
	Code        string `json:"code" example:"NET30"`
	Name        string `json:"name" example:"Net 30"`
	Description string `json:"description" example:"This means the buyer has 30 days from the invoice date to pay"`
}

type CreatePaymentMethodReq struct {
	Code        string `json:"code" example:"BT" validate:"required"`
	Name        string `json:"name" example:"Bank Transfer" validate:"required"`
	Description string `json:"description" example:"Bank Transfer via a local bank"`
}

type UpdatePaymentMethodReq struct {
	Code        string `json:"code" example:"BT"`
	Name        string `json:"name" example:"Bank Transfer"`
	Description string `json:"description" example:"Bank Transfer via a local bank"`
}
