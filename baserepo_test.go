package zeptobaserepo

import (
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

func TestBaseModeRepo_Create(t *testing.T) {
	rider := Rider{
		Name:  "Sajid",
		Phone: "+91-9939879451",
	}
	err := riderRepo.Create(&rider)
	if err != nil {
		t.Fatalf("Create failed")
	} else {
		t.Logf("Create Successful")
	}

}

func TestBaseModeRepo_BatchCreate(t *testing.T) {
	var riders = make([]Rider, 10, 10)
	for i := 0; i < 100; i++ {
		riders = append(riders, Rider{
			Name:  "Sajid",
			Phone: "+91-9939879451",
		})
	}
	t.Logf("len is %d", len(riders))
	err := riderRepo.BatchCreate(riders)
	if err != nil {
		t.Fatalf("Create failed")
	} else {
		t.Logf("Create Successful")
	}
}

func TestBaseModeRepo_FindById(t *testing.T) {

	rider := Rider{
		Name:  "Sajid",
		Phone: "+91-9939879451",
	}
	riderRepo.Create(&rider)
	riderModel, err := riderRepo.FindById(rider.ID)
	if err != nil {
		t.Fatalf("Find by Id failed")
	} else {
		t.Logf("Find by ID success full %v", riderModel)
	}
}

func TestBaseModeRepo_FindAll(t *testing.T) {
	result, err := riderRepo.FindAll()
	if err != nil {
		t.Fatalf("Find all failed")
	} else {
		t.Logf("Find by ID success full %v", len(result.([]Rider)))
	}
}

func TestBaseModeRepo_Update(t *testing.T) {
	rider := Rider{
		Name:  "Sajid",
		Phone: "+91-9939879451",
	}
	riderRepo.Create(&rider)

	rider.Name = "UpdatedName"
	riderRepo.Update(&rider)

	res, _ := riderRepo.FindById(rider.ID)
	if res.(*Rider).Name != "UpdatedName" {
		t.Fatalf("Updated Failed")
	}

}

func TestBaseModeRepo_UpdateSpecificFields(t *testing.T) {

}

func TestBaseModeRepo_Delete(t *testing.T) {

}

func TestBaseModeRepo_DeleteAll(t *testing.T) {

}
