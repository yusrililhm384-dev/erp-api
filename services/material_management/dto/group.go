package dto

type CreateMaterialGroupReq struct {
	Name        string `json:"name" example:"Electronics" validate:"required"`
	Description string `json:"desciption" example:"Electronics and electronical components" validate:"required"`
}

type UpdateMaterialGroupReq struct {
	Name        string `json:"name" example:"Electronics"`
	Description string `json:"desciption" example:"Electronics and electronical components"`
}
