package entity

import (
	"enterprise_resource_planning/services/sales_and_distribution/entity"

	org "enterprise_resource_planning/services/organization/entity"

	"gorm.io/gorm"
)

type MaterialMaster struct {
	gorm.Model

	// basic
	Code        string  `json:"code" gorm:"not null;uniqueIndex"`
	Name        string  `json:"name" gorm:"not null;type:varchar(100)"`
	Description string  `json:"description" gorm:"not null;type:text"`
	Barcode     *string `json:"barcode" gorm:"type:varchar(20)"`

	MaterialTypeID uint          `json:"material_type_id" gorm:"index;not null"`
	MaterialType   *MaterialType `json:"material_type" gorm:"foreignKey:MaterialTypeID"`

	MaterialGroupID uint           `json:"material_group_id" gorm:"index;not null"`
	MaterialGroup   *MaterialGroup `json:"material_group" gorm:"foreignKey:MaterialGroupID"`

	// manufacturer
	ManufacturerID         *uint         `json:"manufacturer_id" gorm:"index"`
	ManufacturerPartNumber *string       `json:"manufacturer_part_number" gorm:"type:varchar(100)"`
	Manufacturer           *Manufacturer `json:"manufacturer" gorm:"foreignKey:ManufacturerID"`

	// brand
	BrandID *uint  `json:"brand_id" gorm:"index"`
	Brand   *Brand `json:"brand" gorm:"foreignKey:BrandID"`

	BaseUnitID uint          `json:"base_unit_id" gorm:"index;not null"`
	BaseUnit   *MaterialUnit `json:"base_unit" gorm:"foreignKey:BaseUnitID"`

	CountryID *uint        `json:"country_code" gorm:"index"`
	Country   *org.Country `json:"country" gorm:"foreignKey:CountryID"`

	MaterialStatusID uint    `json:"material_status_id" gorm:"index;not null"`
	MaterialStatus   *Status `json:"material_status" gorm:"foreignKey:MaterialStatusID"`

	MaterialPhysical   *MaterialPhysical   `json:"material_physical"`
	MaterialPurchasing *MaterialPurchasing `json:"material_purchasing"`

	MaterialConversions []*MaterialConversion `json:"material_conversions"`
	MaterialWarehouses  []*MaterialWarehouse  `json:"material_warehouses"`
}

type MaterialPhysical struct {
	gorm.Model

	MaterialMasterID uint            `json:"material_id" gorm:"uniqueIndex;not null"`
	MaterialMaster   *MaterialMaster `json:"material" gorm:"foreignKey:MaterialMasterID"`

	WeightUnitID *uint         `json:"weight_unit_id" gorm:"index"`
	WeightUnit   *MaterialUnit `json:"weight_unit" gorm:"foreignKey:WeightUnitID"`

	GrossWeight *float64 `json:"gross_weight" gorm:"type:double"`
	NetWeight   *float64 `json:"net_weight" gorm:"type:double"`

	VolumeUnitID *uint         `json:"volume_id" gorm:"index"`
	VolumeUnit   *MaterialUnit `json:"volume_unit" gorm:"foreignKey:VolumeUnitID"`
	Volume       *float64      `json:"volume" gorm:"type:double"`

	DimensionUnitID *uint         `json:"dimension_unit_id" gorm:"index"`
	DimensionUnit   *MaterialUnit `json:"dimension_unit" gorm:"foreignKey:DimensionUnitID"`

	Length *float64 `json:"length" gorm:"type:double"`
	Width  *float64 `json:"width" gorm:"type:double"`
	Height *float64 `json:"height" gorm:"type:double"`

	IsSerialManaged bool `json:"is_serial_managed" gorm:"type:boolean"`
	IsBatchManaged  bool `json:"is_batch_managed" gorm:"type:boolean"`
}

type MaterialConversion struct {
	gorm.Model

	MaterialMasterID uint            `json:"material_master_id" gorm:"uniqueIndex:uq_idx_material_coversions_unit;not null"`
	MaterialMaster   *MaterialMaster `json:"material_master" gorm:"foreignKey:MaterialMasterID"`

	MaterialUnitID uint          `json:"material_unit_id" gorm:"uq_idx_material_coversions_unit;not null"`
	MaterialUnit   *MaterialUnit `json:"material_unit" gorm:"foreignKey:MaterialUnitID"`

	Numerator   uint `json:"numerator" gorm:"type:int;not null;check:,numerator > 0"`
	Denomirator uint `json:"denomirator" gorm:"type:int;not null;check:,denomirator > 0"`
}

type MaterialPurchasing struct {
	gorm.Model

	MaterialMasterID uint            `json:"material_master_id" gorm:"uniqueIndex;not null"`
	MaterialMaster   *MaterialMaster `json:"material_master" gorm:"foreignKey:MaterialMasterID"`

	PurchasingUnitID uint          `json:"purchasing_unit_id" gorm:"index;not null"`
	PurchasingUnit   *MaterialUnit `json:"purchasing_unit" gorm:"foreignKey:PurchasingUnitID"`

	PurchasingOrganizationID uint                    `json:"purchasing_organization_id" gorm:"index;not null"`
	PurchasingOrganization   *PurchasingOrganization `json:"purchasing_organization" gorm:"foreignKey:PurchasingOrganizationID"`

	PurchasingGroupID uint             `json:"purchasing_group_id" gorm:"index;not null"`
	PurchasingGroup   *PurchasingGroup `json:"purchasing_group" gorm:"foreignKey:PurchasingGroupID"`

	OrderUnitID uint          `json:"order_unit_id" gorm:"index;not null"`
	OrderUnit   *MaterialUnit `json:"order_unit" gorm:"foreignKey:OrderUnitID"`

	POText string `json:"po_text" gorm:"type:text;not null"`

	MinimumOrderQty  float64 `json:"minimum_order_quantity" gorm:"type:decimal;not null"`
	OrderMultipleQty float64 `json:"order_multiple_quantity" gorm:"type:decimal;not null"`

	GoodReceiptProcessingTime uint `json:"good_receipt_process_time" gorm:"type:int;not null"`
	PlannedDeliveryDays       uint `json:"planned_delivery_days" gorm:"type:int;not null"`

	IsAutomaticPO bool `json:"is_automatic_po" gorm:"type:boolean;default:false"`

	OverDeliveryTolerance  *float64 `json:"over_delivery_tolerance" gorm:"decimal"`
	UnderDeliveryTolerance *float64 `json:"under_delivery_tolerance" gorm:"decimal"`
}

type MaterialWarehouse struct {
	gorm.Model

	MaterialMasterID uint            `json:"material_master_id" gorm:"index;not null"`
	MaterialMaster   *MaterialMaster `json:"material_master"`

	WarehouseID uint       `json:"warehouse_id" gorm:"index;not null"`
	Warehouse   *Warehouse `json:"warehouse" gorm:"foreignKey:WarehouseID"`

	StorageLocID    uint             `json:"storage_location_id" gorm:"index;not null"`
	StorageLocation *StorageLocation `json:"storage_location" gorm:"foreignKey:StorageLocID"`
}

type MaterialSales struct {
	gorm.Model

	MaterialMasterID uint            `json:"material_master_id" gorm:"index"`
	MaterialMaster   *MaterialMaster `json:"material_master" gorm:"foreignKey:MaterialMasterID"`

	OrderUnitID uint          `json:"order_unit_id" gorm:"index"`
	OrderUnit   *MaterialUnit `json:"order_unit" gorm:"foreignKey:OrderUnitID"`

	DivisionID uint             `json:"division_id" gorm:"index"`
	Division   *entity.Division `json:"division" gorm:"foreignKey:DivisionID"`

	SalesOrgID uint             `json:"sales_org_id" gorm:"index;not null"`
	SalesOrg   *entity.SalesOrg `json:"sales_org" gorm:"foreignKey:SalesOrgID"`

	DistChannelID uint                        `json:"dist_channel_id" gorm:"index;not null"`
	DistChannel   *entity.DistributionChannel `json:"dist_channel" gorm:"foreignKey:DistChannelID"`

	DeliveryPlantID uint       `json:"delivery_plant_id" gorm:"index;not null"`
	DeliveryPlant   *Warehouse `json:"delivery_plant" gorm:"foreignKey:DeliveryPlantID"`
}

type MaterialVendor struct {
	gorm.Model

	MaterialMasterID uint            `json:"material_master_id" gorm:"index"`
	MaterialMaster   *MaterialMaster `json:"material_master" gorm:"foreignKey:MaterialMasterID"`
}
