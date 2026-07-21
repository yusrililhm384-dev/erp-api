package dto

type CreateMaterialTypeReq struct {
	Name        string `json:"name" example:"Raw Material" validate:"required"`
	Description string `json:"description" example:"Raw material is a crude commodity or basic substance used by a company to manufacture finished goods" validate:"required"`
}

type UpdateMaterialTypeReq struct {
	Name        string `json:"name" example:"Raw Material"`
	Description string `json:"description" example:"Raw material is a crude commodity or basic substance used by a company to manufacture finished goods"`
}
