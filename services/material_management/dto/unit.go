package dto

type CreateMaterialUnitReq struct {
	CategoryId  uint   `json:"category_id" example:"1" validate:"required"`
	Code        string `json:"code" example:"EA" validate:"required"`
	Name        string `json:"name" example:"Each" validate:"required"`
	Description string `json:"desciption" example:"Each product" validate:"required"`
}

type MaterialUnitReq struct {
	Code        string `json:"code" example:"EA" validate:"required"`
	Name        string `json:"name" example:"Each" validate:"required"`
	Description string `json:"desciption" example:"Each product" validate:"required"`
}

type UpdateMaterialUnitReq struct {
	CategoryId  uint   `json:"category_id" example:"1"`
	Code        string `json:"code" example:"EA"`
	Name        string `json:"name" example:"Each"`
	Description string `json:"desciption" example:"Each product"`
}

type CreateMaterialUnitCategoryReq struct {
	Code          string             `json:"code" example:"W" validate:"required"`
	Name          string             `json:"name" example:"Weight" validate:"required"`
	MaterialUnits []*MaterialUnitReq `json:"material_units"`
}

type UpdateMaterialUnitCategoryReq struct {
	Code string `json:"code" example:"W"`
	Name string `json:"name" example:"Weight"`
}
