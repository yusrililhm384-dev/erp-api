package dto

import "time"

type CreateMutationReq struct {
	BranchID     uint       `json:"branch_id" example:"1" validate:"required"`
	PositionID   uint       `json:"position_id" example:"1" validate:"required"`
	TypeID       uint       `json:"type_id" example:"1" validate:"required"`
	DepartmentID uint       `json:"department_id" example:"1" validate:"required"`
	StartDate    time.Time  `json:"start_date" example:"2026-01-01T00:00:00Z" validate:"required"`
	EndDate      *time.Time `json:"end_date" example:"2026-01-01T00:00:00Z"`
}
