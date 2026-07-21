package dto

type CreateStatusReq struct {
	Code string `json:"code" example:"ACTIVE" validate:"required"`
	Name string `json:"name" example:"Active" validate:"required"`
}

type UpdateStatusReq struct {
	Code string `json:"code" example:"ACTIVE"`
	Name string `json:"name" example:"Active"`
}
