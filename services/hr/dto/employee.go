package dto

import "time"

type CreateEmployeeReq struct {
	Name         string             `json:"name" example:"Xiaoching" validate:"required"`
	Email        string             `json:"email" example:"xiao@mail.com" validate:"email,required"`
	Phone        string             `json:"phone" example:"08123456789" validate:"required"`
	Gender       string             `json:"gender" example:"male" validate:"required"`
	Address      string             `json:"address" example:"South Jakarta" validate:"required"`
	BirthDate    time.Time          `json:"birth_date" example:"2001-01-01T00:00:00Z" validate:"required"`
	JoinDate     time.Time          `json:"join_date" example:"2026-01-01T00:00:00Z" validate:"required"`
	TempUsername string             `json:"temporary_username" example:"tempuser" validate:"required"`
	TempPassword string             `json:"temporary_password" example:"temppassword" validate:"required"`
	Legal        *CreateLegalReq    `json:"legal" validate:"required"`
	Mutation     *CreateMutationReq `json:"mutation" validate:"required"`
}

type UpdateEmployeeReq struct {
	Name      string    `json:"name,omitempty" example:"Xiaoching"`
	Email     string    `json:"email,omitempty" example:"xiao@mail.com"`
	Phone     string    `json:"phone,omitempty" example:"08123456789"`
	Gender    string    `json:"gender" example:"male"`
	Address   string    `json:"address,omitempty" example:"South Jakarta"`
	BirthDate time.Time `json:"birth_date" example:"2001-01-01T00:00:00Z"`
	JoinDate  time.Time `json:"join_date" example:"2026-01-01T00:00:00Z"`
}
