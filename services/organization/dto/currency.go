package dto

type CreateCurrencyReq struct {
	Code      string `json:"code" example:"IDR" validate:"required"`
	Name      string `json:"name" example:"Indonesia" validate:"required"`
	CountryId uint   `json:"country_id" example:"1" validate:"required"`
}

type UpdateCurrencyReq struct {
	Code      string `json:"code" example:"IDR"`
	Name      string `json:"name" example:"Indonesia Rupiah"`
	CountryId uint   `json:"country_id" example:"1"`
}
