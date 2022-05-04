package zeptobaserepo

import (
	"encoding/json"
	"gorm.io/gorm"
)

type RiderEnv struct {
	gorm.Model
}

func init() {

}

func (rider *Rider) AfterCreate(tx *gorm.DB) error {

	// publish events
	jsonByte, _ := json.MarshalIndent(rider, "", "")
	topicStr := "rider-events"
	rider.PublishEvents(jsonByte, topicStr)

	// save versions
	riderEnv := RiderEnv{}
	return riderEnvRepo.Create(&riderEnv)
}
