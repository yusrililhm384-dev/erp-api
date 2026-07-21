package dto

type CreateMaterialMasterReq struct {
	Code                   string                     `json:"code" example:"MBA-M4-16-256-MIDNIGHT" validate:"required"`
	Name                   string                     `json:"name" example:"Macbook Air M4 16GB/256BG MIDNIGHT" validate:"required"`
	Description            string                     `json:"description" example:"Macbook Air with M4 chip and 16 GB Memory" validate:"required"`
	Barcode                *string                    `json:"barcode" example:"1111222244445"`
	MaterialTypeId         uint                       `json:"material_type_id" example:"1" validate:"required"`
	MaterialGroupId        uint                       `json:"material_group_id" example:"1"`
	ManufacturerId         *uint                      `json:"manufacturer_id" example:"1"`
	ManufacturerPartNumber *string                    `json:"manufacturer_part_number" example:"xxxxxxxxxxxx"`
	BrandId                *uint                      `json:"brand_id" example:"1"`
	BaseUnitId             uint                       `json:"base_unit_id" example:"1" validate:"required"`
	CountryId              *uint                      `json:"country_id" example:"1"`
	MaterialStatusId       uint                       `json:"material_status_id" example:"1" validate:"required"`
	MaterialPhysical       *CreateMaterialPhysicalReq `json:"material_physical,omitempty"`
}

type UpdateMaterialMasterReq struct {
	Code                   string  `json:"code" example:"MBA-M4-16-256-MIDNIGHT"`
	Name                   string  `json:"name" example:"Macbook Air M4 16GB/256BG MIDNIGHT"`
	Description            string  `json:"description" example:"Macbook Air with M4 chip and 16 GB Memory"`
	Barcode                *string `json:"barcode" example:"1111222244445"`
	MaterialTypeId         uint    `json:"material_type_id" example:"1"`
	MaterialGroupId        uint    `json:"material_group_id" example:"1"`
	ManufacturerId         *uint   `json:"manufacturer_id" example:"1"`
	ManufacturerPartNumber *string `json:"manufacturer_part_number" example:"xxxxxxxxxxxx"`
	BrandId                *uint   `json:"brand_id" example:"1"`
	BaseUnitId             uint    `json:"base_unit_id" example:"1"`
	CountryId              *uint   `json:"country_id" example:"1"`
	MaterialStatusId       uint    `json:"material_status_id" example:"1"`
}

type CreateMaterialPhysicalReq struct {
	WeightUnitId    *uint    `json:"weight_unit_id,omitempty" example:"1" validate:"required"`
	GrossWeight     *float64 `json:"gross_weight" example:"5" validate:"required"`
	NetWeight       *float64 `json:"net_weight" example:"4.5" validate:"required"`
	VolumeUnitId    *uint    `json:"volume_unit_id" example:"1" validate:"required"`
	Volume          *float64 `json:"volume" example:"1" validate:"required"`
	DimensionUnitId *uint    `json:"dimension_unit_id" example:"1" validate:"required"`
	Width           *float64 `json:"width" example:"1" validate:"required"`
	Length          *float64 `json:"length" example:"1" validate:"required"`
	Height          *float64 `json:"height" example:"1" validate:"required"`
	IsSerialManaged bool     `json:"is_serial_managed" example:"false" validate:"required"`
	IsBatchManaged  bool     `json:"is_batch_managed" example:"false" validate:"required"`
}

type UpdateMaterialPhysicalReq struct {
	WeightUnitId    *uint    `json:"weight_unit_id,omitempty" example:"1"`
	GrossWeight     *float64 `json:"gross_weight" example:"5"`
	NetWeight       *float64 `json:"net_weight" example:"4.5"`
	VolumeUnitId    *uint    `json:"volume_unit_id" example:"1"`
	Volume          *float64 `json:"volume" example:"1"`
	DimensionUnitId *uint    `json:"dimension_unit_id" example:"1"`
	Width           *float64 `json:"width" example:"1"`
	Length          *float64 `json:"length" example:"1"`
	Height          *float64 `json:"height" example:"1"`
	IsSerialManaged bool     `json:"is_serial_managed" example:"false"`
	IsBatchManaged  bool     `json:"is_batch_managed" example:"false"`
}

type CreateMaterialConversionReq struct {
	MaterialUnitId uint `json:"material_unit_id" example:"1" validate:"required"`
	Numerator      uint `json:"numerator" example:"10" validate:"gt=0,required"`
	Denomirator    uint `json:"denomirator" example:"1" validate:"gt=0,required"`
}

type UpdateMaterialConversionReq struct {
	MaterialUnitId uint `json:"material_unit_id" example:"1"`
	Numerator      uint `json:"numerator" example:"10" validate:"gt=0"`
	Denomirator    uint `json:"denomirator" example:"1" validate:"gt=0"`
}
