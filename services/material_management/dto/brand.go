package dto

type BrandReq struct {
	Name string `json:"name" example:"Apple" validate:"required"`
}
