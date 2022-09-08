package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Identity  uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
