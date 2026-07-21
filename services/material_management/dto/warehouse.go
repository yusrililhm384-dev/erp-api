package dto

type CreateWarehouseReq struct {
	Name        string `json:"name" validate:"required" example:"Main"`
	Description string `json:"desciption" validate:"required" example:"Main Warehouse"`
	Address     string `json:"address" validate:"required" example:"Jl. Silang Monas, Gambir, Kecamatan Gambir, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10110, Indonesia"`
	Location    *Point `json:"location" validate:"required"`
}

type UpdateWarehouseReq struct {
	Name        string `json:"name" example:"Main"`
	Description string `json:"desciption" example:"Main Warehouse"`
	Address     string `json:"address" example:"Jl. Silang Monas, Gambir, Kecamatan Gambir, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10110, Indonesia"`
	Location    *Point `json:"location"`
}
