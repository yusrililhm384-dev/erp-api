package dto

type CreateCountryReq struct {
	Code string `json:"code" example:"ID" validate:"required"`
	Name string `json:"name" example:"Indonesia" validate:"required"`
}

type UpdateCountryReq struct {
	Code string `json:"code" example:"ID"`
	Name string `json:"name" example:"Indonesia"`
}
