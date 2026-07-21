package dto

type CreateLanguageReq struct {
	Code string `json:"code" example:"INDO" validate:"required"`
	Name string `json:"name" example:"Indonesia" validate:"required"`
}

type UpdateLanguageReq struct {
	Code string `json:"code" example:"INDO"`
	Name string `json:"name" example:"Indonesia"`
}
