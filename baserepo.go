package zeptorepocommons

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
)

type CrudRepo interface {
	Create(any)
	BatchCreate(any)
	FindById(uint) any
	FindAll() (*PaginatorQueryResult, error)
	Query(condition *QueryCondition) (*PaginatorQueryResult, error)
	Update()
	UpdateSpecificFields(uint, map[string]any)
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

func GetRepo(db *gorm.DB, typ reflect.Type) *BaseRepo {
	DefaultBatchCreateSize = 100
	DefaultPageSize = 100
	return &BaseRepo{db: db, baseModelType: typ}
}

func (bmr BaseRepo) Create(value any) error {
	result := bmr.db.Create(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr BaseRepo) BatchCreate(value any) error {
	result := bmr.db.CreateInBatches(value, DefaultBatchCreateSize)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr BaseRepo) FindById(id uint) (any, error) {
	value := reflect.New(bmr.baseModelType).Interface()
	result := bmr.db.
		Preload(clause.Associations).
		Find(value, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return value, nil
}

func (bmr BaseRepo) FindAll(offset int) (*PaginatorQueryResult, error) {

	values := reflect.New(reflect.SliceOf(bmr.baseModelType)).Interface()
	nextOffset := offset + DefaultPageSize + 1
	result := bmr.db.
		Preload(clause.Associations).
		Offset(offset).Limit(DefaultPageSize).
		Find(values)
	if result.Error != nil {
		return nil, result.Error
	}
	return &PaginatorQueryResult{values, nextOffset}, nil
}

func (bmr BaseRepo) Query(query *Query) (*PaginatorQueryResult, error) {
	var (
		offset   = 0
		pageSize = DefaultPageSize
	)
	if query.pageConfig != nil {
		offset = query.pageConfig.offset
		pageSize = query.pageConfig.limit
	}
	values := reflect.New(reflect.SliceOf(bmr.baseModelType)).Interface()

	db := bmr.db.
		Preload(clause.Associations).
		Where(query.queryCondition.getPreparedStatement()).
		Offset(offset).
		Limit(pageSize)

	if query.pageConfig.orderBy != nil {
		for k, v := range query.pageConfig.orderBy {
			d := false
			if strings.EqualFold(v, "desc") {
				d = true
			}
			db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: k}, Desc: d})
		}
	}

	result := db.Find(values)
	if result.Error != nil {
		return nil, result.Error
	}
	nextOffset := offset + pageSize + 1
	return &PaginatorQueryResult{values, nextOffset}, nil
}

func (bmr BaseRepo) Update(value interface{}) error {

	result := bmr.db.Updates(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr BaseRepo) UpdateSpecificFields(id uint, fields map[string]interface{}) error {
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

func (bmr BaseRepo) Delete(value interface{}) error {
	result := bmr.db.Delete(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bmr BaseRepo) DeleteALl() error {
	value := reflect.New(bmr.baseModelType).Interface()
	result := bmr.db.Where("1 = 1").Delete(value)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
