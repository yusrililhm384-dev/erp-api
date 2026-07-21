package dto

type CreateBranchReq struct {
	CompanyId   uint   `json:"company_id" example:"1" validate:"required"`
	Name        string `json:"name" validate:"required" example:"Main"`
	Description string `json:"desciption" validate:"required" example:"Main Branch"`
	Address     string `json:"address" validate:"required" example:"Jl. Silang Monas, Gambir, Kecamatan Gambir, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10110, Indonesia"`
	Location    *Point `json:"location" validate:"required"`
}

type UpdateBranchReq struct {
	CompanyId   uint   `json:"company_id" example:"1"`
	Name        string `json:"name" example:"Main"`
	Description string `json:"desciption" example:"Main Branch"`
	Address     string `json:"address" example:"Jl. Silang Monas, Gambir, Kecamatan Gambir, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10110, Indonesia"`
	Location    *Point `json:"location"`
}
