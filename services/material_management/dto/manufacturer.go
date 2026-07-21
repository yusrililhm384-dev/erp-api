package dto

type CreateManufacturerReq struct {
	Code string `json:"code" example:"Apple" validate:"required"`
	Name string `json:"name" example:"Apple" validate:"required"`
}

type UpdateManufacturerReq struct {
	Code string `json:"code" example:"Apple"`
	Name string `json:"name" example:"Apple"`
}
