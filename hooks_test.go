package zeptorepocommons

import (
	"encoding/json"
	"gorm.io/gorm"
)

type RiderEnv struct {
	PID   uint  `gorm:"primarykey"`
	Rider Rider `gorm:"embedded;embeddedPrefix:rider_"`
}

func init() {

}

func (rider *Rider) AfterCreate(tx *gorm.DB) error {

	//publish events
	jsonByte, _ := json.MarshalIndent(rider, "", "")
	rider.PublishEvents(jsonByte)

	// save versions
	riderEnv := RiderEnv{}
	dto.Map(&riderEnv, *rider)
	return riderEnvRepo.Create(&riderEnv)
}
