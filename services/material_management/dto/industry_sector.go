package dto

type CreateIndustrySectorReq struct {
	Name        string `json:"name" example:"OGE" validate:"required"`
	Description string `json:"description" example:"OGE is oil, gas, and energy industrial sector." validate:"required"`
}

type UpdateIndustrySectorReq struct {
	Name        string `json:"name" example:"OGE"`
	Description string `json:"description" example:"OGE is oil, gas, and energy industrial sector."`
}
