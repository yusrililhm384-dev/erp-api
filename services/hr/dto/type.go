package dto

import "time"

type CreateTypeReq struct {
	Name        string `json:"name" example:"Full Time" validate:"required"`
	Description string `json:"description" example:"Full time employee" validate:"required"`
}

type UpdateTypeReq struct {
	Name        string `json:"name" example:"Full Time"`
	Description string `json:"description" example:"Full time employee"`
}

type DetailTypeRes struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
