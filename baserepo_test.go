package zeptobaserepo

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

type Rider struct {
	gorm.Model
	Name  string
	Phone string
}

type RiderRepo struct {
	*BaseRepo
}

var (
	riderRepo *RiderRepo
)

func init() {
	dsn := "host=localhost user=postgres password=password dbname=zeptodb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	riderRepo = &RiderRepo{getRepo(db, reflect.TypeOf(Rider{}))}

}

func getSampleRider() *Rider {
	return &Rider{
		Name:  "Sajid",
		Phone: "+91-9939879451",
	}
}

func TestBaseModeRepo_Create(t *testing.T) {
	rider := getSampleRider()
	err := riderRepo.Create(rider)
	assert.Nil(t, err)
}

func TestBaseModeRepo_BatchCreate(t *testing.T) {
	var riders = make([]Rider, 10, 10)
	for i := 0; i < 100; i++ {
		riders = append(riders, *getSampleRider())
	}
	err := riderRepo.BatchCreate(riders)
	assert.Nil(t, err)
}

func TestBaseModeRepo_FindById(t *testing.T) {
	rider := getSampleRider()
	riderRepo.Create(rider)
	riderModel, err := riderRepo.FindById(rider.ID)
	assert.Nil(t, err)
	assert.Equal(t, riderModel.(*Rider).ID, rider.ID)
}

func TestBaseModeRepo_FindAll(t *testing.T) {
	paginatedResult, err := riderRepo.FindAll(100)
	if err != nil {
		t.Fatalf("Find all failed")
	} else {
		assert.Equal(t, len(paginatedResult.values.([]Rider)), DefaultPageSize)
		assert.Equal(t, paginatedResult.nextOffset, 201)
	}
}

//func TestBaseModeRepo_Update(t *testing.T) {
//	rider := Rider{
//		Name:  "Sajid",
//		Phone: "+91-9939879451",
//	}
//	riderRepo.Create(&rider)
//
//	rider.Name = "UpdatedName"
//	riderRepo.Update(&rider)
//
//	res, _ := riderRepo.FindById(rider.ID)
//	if res.(*Rider).Name != "UpdatedName" {
//		t.Fatalf("Updated Failed")
//	}
//
//}

func TestBaseModeRepo_UpdateSpecificFields(t *testing.T) {

}

func TestBaseModeRepo_Delete(t *testing.T) {

}

func TestBaseModeRepo_DeleteAll(t *testing.T) {

}
