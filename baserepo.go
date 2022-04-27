package zeptobaserepo

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
)

type CrudRepo interface {
	Create(interface{})
	BatchCreate(interface{})

	FindById(uint) interface{}
	FindAll() (*PaginatorQueryResult, error)
	Query(condition SearchCondition) (*PaginatorQueryResult, error)

	Update()
	UpdateSpecificFields(uint, map[string]interface{})

	Delete(uint)
	DeleteAll()
}

var (
	DefaultBatchCreateSize int
	DefaultPageSize        int
)

type BaseRepo struct {
	baseModelType reflect.Type
	db            *gorm.DB
}

func getRepo(db *gorm.DB, typ reflect.Type) *BaseRepo {
	DefaultBatchCreateSize = 100
	DefaultPageSize = 100
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
	result := bmr.db.CreateInBatches(value, DefaultBatchCreateSize)
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

// FindAll finds all the instances of the model.
// The first and only argument is the offset
func (bmr *BaseRepo) FindAll(params ...int) (*PaginatorQueryResult, error) {
	values := reflect.MakeSlice(reflect.SliceOf(bmr.baseModelType), 0, 0).Interface()
	var offset, nextOffset int
	switch len(params) {
	case 0:
		offset = 0
	case 1:
		offset = params[0]
	default:
		return nil, errors.New("specified more number of arguments than specified")
	}
	nextOffset = offset + DefaultPageSize + 1
	result := bmr.db.Offset(offset).Limit(DefaultPageSize).Find(&values)
	if result.Error != nil {
		return nil, result.Error
	}
	return &PaginatorQueryResult{values, nextOffset}, nil
}

func (bmr *BaseRepo) Query(condition SearchCondition) (*PaginatorQueryResult, error) {
	values := reflect.MakeSlice(reflect.SliceOf(bmr.baseModelType), 0, 0).Interface()
	result := bmr.db.Where(condition.getPreparedStatement()).Offset(condition.offset).Limit(DefaultPageSize).Find(&values)
	if result.Error != nil {
		return nil, result.Error
	}
	nextOffset := condition.offset + DefaultPageSize + 1
	return &PaginatorQueryResult{values, nextOffset}, nil
}

func (bmr *BaseRepo) Update(value interface{}) error {
	result := bmr.db.Updates(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr *BaseRepo) UpdateSpecificFields(id uint, fields map[string]interface{}) error {

	// do this in continuous session mode
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
