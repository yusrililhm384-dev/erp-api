package entity

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Action    string         `json:"action" gorm:"type:varchar(25);not null"`
	TableName string         `json:"table_name" gorm:"type:varchar(50);not null"`
	RecordId  uint           `json:"record_id" gorm:"not null"`
	OldValue  datatypes.JSON `json:"old_value" gorm:"type:jsonb"`
	NewValue  datatypes.JSON `json:"new_value" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at"`
}
