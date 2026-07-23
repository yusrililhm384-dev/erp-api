package dto

type CreateVendorReq struct {
	Code          string              `json:"code" example:"VMETRODATA" validate:"required"`
	Name          string              `json:"name" example:"PT Xiaoching Mineral" validate:"required"`
	Description   string              `json:"description" example:"Coal supplier" validate:"required"`
	TaxNumber     *string             `json:"tax_number" example:"xxxxaaaabbbbcccc"`
	Address       string              `json:"address" example:"Jl. Silang Monas, Gambir, Kecamatan Gambir, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10110, Indonesia" validate:"required"`
	Location      *Point              `json:"location" validate:"required"`
	VendorTypeId  uint                `json:"vendor_type_id" example:"1" validate:"required"`
	CountryId     uint                `json:"country_id" example:"1" validate:"required"`
	LanguageId    uint                `json:"language_id" example:"1" validate:"required"`
	CurrencyId    uint                `json:"currency_id" example:"1" validate:"required"`
	PaymentTermId uint                `json:"payment_term_id" example:"1" validate:"required"`
	Contacts      []*ContactVendorReq `json:"contacts" validate:"requied"`
}

type UpdateVendorReq struct {
	Code          string  `json:"code" example:"VMETRODATA"`
	Name          string  `json:"name" example:"PT Xiaoching Mineral"`
	Description   string  `json:"description" example:"Coal supplier"`
	TaxNumber     *string `json:"tax_number" example:"xxxxaaaabbbbcccc"`
	Address       string  `json:"address" example:"Jl. Silang Monas, Gambir, Kecamatan Gambir, Kota Jakarta Pusat, Daerah Khusus Ibukota Jakarta 10110, Indonesia"`
	Location      *Point  `json:"location"`
	VendorTypeId  uint    `json:"vendor_type_id" example:"1"`
	CountryId     uint    `json:"country_id" example:"1"`
	LanguageId    uint    `json:"language_id" example:"1"`
	CurrencyId    uint    `json:"currency_id" example:"1"`
	PaymentTermId uint    `json:"payment_term_id" example:"1"`
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
