package zeptobaserepo

import (
	"gorm.io/gorm"
	"reflect"
)

type CrudRepo interface {
	Create(interface{})
	BatchCreate(interface{})
	FindById(uint) interface{}
	FindAll() interface{}
	Update()
	UpdateSpecificFields(uint, map[string]interface{})
	Delete(uint)
	DeleteAll()
}

var (
	CreateBatchSize int
)

type BaseRepo struct {
	baseModelType reflect.Type
	db            *gorm.DB
}

func getRepo(db *gorm.DB, typ reflect.Type) *BaseRepo {
	CreateBatchSize = 100
	return &BaseRepo{db: db, baseModelType: typ}
}

func (bmr *BaseRepo) Create(value interface{}) error {
	result := bmr.db.Create(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr *BaseRepo) BatchCreate(value interface{}) error {
	result := bmr.db.CreateInBatches(value, CreateBatchSize)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr *BaseRepo) FindById(id uint) (interface{}, error) {
	value := reflect.New(bmr.baseModelType).Interface()
	result := bmr.db.Find(&value, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}

func (bmr *BaseRepo) FindAll() (interface{}, error) {
	values := reflect.MakeSlice(reflect.SliceOf(bmr.baseModelType), 10, 10).Interface()
	result := bmr.db.Find(&values)
	if result.Error != nil {
		return nil, result.Error
	}
	return values, nil
}

func (bmr *BaseRepo) Update(value interface{}) error {
	result := bmr.db.Updates(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr *BaseRepo) UpdateSpecificFields(id uint, fields map[string]interface{}) error {

	res, _ := bmr.FindById(id)
	if res != nil {
		result := bmr.db.Model(res).Updates(fields)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
	return nil
}

func (bmr *BaseRepo) Delete(value interface{}) error {
	result := bmr.db.Delete(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr *BaseRepo) DeleteALl() error {
	value := reflect.New(bmr.baseModelType).Interface()
	result := bmr.db.Where("1 = 1").Delete(&value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
