package zeptorepocommons

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	gorm.Model
	Version uint `gorm:"autoIncrement"`
}

type NonIdBaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
	Version   uint         `gorm:"autoIncrement"`
}
