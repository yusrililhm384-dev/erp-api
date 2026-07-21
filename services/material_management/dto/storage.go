package dto

type CreateStorageLocationReq struct {
	Code        string `json:"code" example:"RM" validate:"required"`
	Name        string `json:"name" example:"Raw Material" validate:"required"`
	Description string `json:"desciption" example:"StoLct. Raw material" validate:"required"`
}

type UpdateStorageLocationReq struct {
	Code        string `json:"code" example:"RM"`
	Name        string `json:"name" example:"Raw Material"`
	Description string `json:"desciption" example:"StoLct. Raw material"`
}
