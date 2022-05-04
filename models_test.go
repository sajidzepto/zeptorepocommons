package zeptobaserepo

import "zeptobaserepo/models"

type Rider struct {
	models.BaseModel
	Name     string
	Phone    string
	VendorId uint
	// belongs to relation
	Vendor *RiderVendor
	// has one relation
	Identification *IdentificationModel
	// has many
	Addresses []AddressModel
	// many to many
	Stores []StoreModel `gorm:"many2many:rider_store_mapping;"`
}

type RiderVendor struct {
	models.BaseModel
	Name  string
	Phone string
}

type IdentificationModel struct {
	models.BaseModel
	IdentificationType string
	IdentificationId   string
	RiderId            uint
}

type AddressModel struct {
	models.BaseModel
	AddressString string
	City          string
	Pincode       string
	RiderId       uint
}

type StoreModel struct {
	models.BaseModel
	Pincode string
}

type RiderRepo struct {
	*BaseRepo
}

type RiderEnvRepo struct {
	*BaseRepo
}

type RiderVendorRepo struct {
	*BaseRepo
}

type IdentificationModelRepo struct {
	*BaseRepo
}

type AddressRepo struct {
	*BaseRepo
}

type StoreModelRepo struct {
	*BaseRepo
}
