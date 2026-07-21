package dto

type CreateDepartmentReq struct {
	Name        string `json:"name" example:"FI" validate:"required"`
	Description string `json:"description" example:"Finance" validate:"required"`
}

type UpdateDepartmentReq struct {
	Name        string `json:"name" example:"FI"`
	Description string `json:"description" example:"Finance"`
}
