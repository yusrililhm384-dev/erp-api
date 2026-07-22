package dto

type CreateVendorReq struct {
	Name        string              `json:"name" example:"PT Xiaoching Mineral" validate:"required"`
	Description string              `json:"description" example:"Coal supplier" validate:"required"`
	Address     string              `json:"address" example:"Jl. Silang Monas, Gambir, Kecamatan Gambir, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10110, Indonesia" validate:"required"`
	Location    *Point              `json:"location" validate:"required"`
	Contacts    []*ContactVendorReq `json:"contacts" validate:"requied"`
}

type UpdateVendorReq struct {
	Name        string `json:"name" example:"PT Xiaoching Mineral"`
	Description string `json:"description" example:"Coal supplier"`
	Address     string `json:"address" example:"Jl. Silang Monas, Gambir, Kecamatan Gambir, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10110, Indonesia"`
	Location    *Point `json:"location"`
}

type ContactVendorReq struct {
	Name  string `json:"name" example:"Xiao Chi" validate:"required"`
	Email string `json:"email" example:"xiao@xmineral.com" validate:"email,required"`
	Phone string `json:"phone" example:"" validate:"required"`
}

type UpdateContactVendorReq struct {
	Name  string `json:"name" example:"Xiao Chi"`
	Email string `json:"email" example:"xiao@xmineral.com"`
	Phone string `json:"phone" example:"081234567890"`
}
