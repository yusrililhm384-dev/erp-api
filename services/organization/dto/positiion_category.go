package dto

type PositionCategoryReq struct {
	Name string `json:"name" example:"Engineer" validate:"required"`
}
