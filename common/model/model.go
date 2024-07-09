package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement;comment:自增主键"`            // 自增主键
	CreatedAt time.Time      `json:"createdAt" gorm:"type:datetime;autoCreateTime;comment:创建时间"` // 创建时间
	UpdatedAt time.Time      `json:"updatedAt" gorm:"type:datetime;autoUpdateTime;comment:更新时间"` // 更新时间
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"type:datetime;index;comment:删除时间"`          // 删除时间
}
