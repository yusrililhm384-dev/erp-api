package dto

type CreatePositionReq struct {
	Name        string `json:"name" example:"Backend Engineer" validate:"required"`
	Description string `json:"description" example:"Create API" validate:"required"`
	ApprovalId  *uint  `json:"approval_id" example:"1" validate:"required"`
	CategoryId  uint   `json:"category_id" example:"1" validate:"required"`
}

type UpdatePositionReq struct {
	Name        string `json:"name" example:"Backend Engineer"`
	Description string `json:"description" example:"Create API"`
	ApprovalId  *uint  `json:"approval_id" example:"1"`
	CategoryId  uint   `json:"category_id" example:"1"`
}
