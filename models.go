package zeptorepocommons

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type IBaseModel interface {
	GetUniqueColumn() *[]clause.Column
}

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
