package zeptobaserepo

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

type Rider struct {
	gorm.Model
	Name     string
	Phone    string
	VendorId uint
	Vendor   RiderVendor
}

type RiderVendor struct {
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
	//db.AutoMigrate(&Rider{}, &RiderVendor{})
	riderRepo = &RiderRepo{getRepo(db, reflect.TypeOf(Rider{}))}

}

func getSampleRider() *Rider {
	return &Rider{
		Name:  "Sajid",
		Phone: "+91-9939879451",
		Vendor: RiderVendor{
			Name:  "Vendor",
			Phone: "VendorPhone",
		},
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func TestBaseRepo_Create(t *testing.T) {
	rider := getSampleRider()
	err := riderRepo.Create(rider)
	assert.Nil(t, err)
}

func TestBaseRepo_BatchCreate(t *testing.T) {
	var riders = make([]Rider, 0, 0)
	for i := 0; i < 100; i++ {
		riders = append(riders, *getSampleRider())
	}
	err := riderRepo.BatchCreate(riders)
	assert.Nil(t, err)
}

func TestBaseRepo_FindById(t *testing.T) {
	rider := getSampleRider()
	riderRepo.Create(rider)

	riderModel, err := riderRepo.FindById(rider.ID)

	assert.Nil(t, err)
	t.Logf("the fetched entity is \n %+v \n", prettyPrint(riderModel))
	assert.Equal(t, riderModel.(*Rider).ID, rider.ID)
}

func TestBaseRepo_FindAll(t *testing.T) {
	paginatedResult, err := riderRepo.FindAll(0)
	assert.Nil(t, err)
	//t.Logf("the fetched entities is \n %+v \n", prettyPrint(paginatedResult.values))
	assert.Equal(t, len(*(paginatedResult.values).(*[]Rider)), DefaultPageSize)
	assert.Equal(t, paginatedResult.nextOffset, 101)
}

func TestBaseRepo_Query(t *testing.T) {

	seachCond := QueryCondition{
		orConditions: []OrConditions{
			{andConditions: []AndConditions{
				{conditions: []Condition{SearchOperatorCondition{
					field:    "Name",
					operator: "=",
					value:    "Sajid",
				}}}}},
		},
	}
	query := Query{
		queryCondition: &seachCond,
		pageConfig: &PageConfig{
			orderBy: map[string]string{
				"id": "DESC",
			},
			offset: 0,
			limit:  100,
		},
	}
	paginatedResult, err := riderRepo.Query(&query)
	assert.Nil(t, err)
	//t.Logf("the fetched entities is \n %+v \n", prettyPrint(paginatedResult.values))
	assert.Equal(t, len(*(paginatedResult.values).(*[]Rider)), 100)
	assert.Equal(t, paginatedResult.nextOffset, 101)
}

func TestBaseRepo_Update(t *testing.T) {
	rider := getSampleRider()
	riderRepo.Create(&rider)
	rider.Name = "UpdatedName"
	riderRepo.Update(rider)
	res, _ := riderRepo.FindById(rider.ID)
	assert.Equal(t, res.(*Rider).Name, "UpdatedName")
}

func TestBaseRepo_UpdateSpecificFields(t *testing.T) {
	rider := getSampleRider()
	riderRepo.Create(rider)
	m := map[string]interface{}{
		"Name": "UpdatedName",
	}
	riderRepo.UpdateSpecificFields(rider.ID, m)
	res, _ := riderRepo.FindById(rider.ID)
	assert.Equal(t, res.(*Rider).Name, "UpdatedName")
}

func TestBaseRepo_Delete(t *testing.T) {
	rider := getSampleRider()
	riderRepo.Create(rider)
	id := rider.ID
	riderRepo.Delete(rider)
	t.Logf("The id the was created %d", id)
	riderModel, err := riderRepo.FindById(id)
	assert.Nil(t, err)
	assert.Nil(t, riderModel)

}

func TestBaseRepo_DeleteAll(t *testing.T) {
	rider := getSampleRider()
	riderRepo.Create(rider)
	riderRepo.DeleteALl()
	pgqr, err := riderRepo.FindAll(0)
	assert.Nil(t, err)
	assert.Equal(t, len(*(pgqr.values).(*[]Rider)), 0)

}
