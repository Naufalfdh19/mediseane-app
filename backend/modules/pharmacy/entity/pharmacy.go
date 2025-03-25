package entity

import (
	"mime/multipart"
	"montelukast/pkg/pagination"

	"github.com/shopspring/decimal"
)

type LogisticParter struct {
	Name string
}

type Pharmacy struct {
	ID            int
	PartnerID     int
	PartnerName   string
	Name          string
	Address       string
	Province      string
	ProvinceID    int
	District      string
	DistrictID    int
	SubDistrict   string
	SubDistrictID int
	City          string
	CityID        int
	Latitude      string
	Longitude     string
	PostalCode    int
	IsActive      bool
	Location      string
	Logo          string
	UpdatedAt     string
}

type PharmacyLogisticPartners struct {
	ID                int
	PharmacyID        int
	LogisticPartnerID int
}

type InternalLogistic struct {
	ID    int
	Name  string
	Price decimal.Decimal
}

type Pagination struct {
	CurrentPage int
	TotalItem   int
	TotalPage   int
}

type PaginatedPharmacies struct {
	Pagination pagination.Pagination
	Pharmacies []Pharmacy
}

type File struct {
	File multipart.File
	ID   int
}
