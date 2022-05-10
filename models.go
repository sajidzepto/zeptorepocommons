package zeptorepocommons

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	Version uint `gorm:"autoIncrement"`
}
