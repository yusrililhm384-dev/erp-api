package entity

import "gorm.io/gorm"

type StorageLocation struct {
	gorm.Model
	Code        string     `json:"code" gorm:"not null;type:varchar(50);uniqueIndex:idx_storage_location_code,where:deleted_at IS NUL"`
	Name        string     `json:"name" gorm:"not null;type:varchar(100)"`
	Description string     `json:"description" gorm:"not null;type:text"`
	WarehouseID uint       `json:"warehouse_id" gorm:"not null;uniqueIndex:idx_storage_location_code,where:deleted_at IS NUL"`
	Warehouse   *Warehouse `json:"warehouse" gorm:"foreignKey:WarehouseID"`
}
