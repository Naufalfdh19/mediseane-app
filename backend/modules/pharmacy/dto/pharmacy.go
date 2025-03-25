package dto

import (
	"mime/multipart"
	"montelukast/pkg/pagination"
)

type Pharmacy struct {
	ID            int    `json:"id"`
	PartnerID     int    `json:"partner_id" binding:"required"`
	Name          string `json:"name" binding:"required,min=3"`
	Address       string `json:"address" binding:"required"`
	Province      string `json:"province" binding:"required"`
	ProvinceID    int    `json:"province_id"  binding:"required"`
	District      string `json:"district" binding:"required"`
	DistrictID    int    `json:"district_id" binding:"required"`
	SubDistrict   string `json:"sub_district" binding:"required"`
	SubDistrictID int    `json:"sub_district_id" binding:"required"`
	City          string `json:"city" binding:"required"`
	CityID        int    `json:"city_id" binding:"required"`
	Latitude      string `json:"latitude" binding:"required"`
	Longitude     string `json:"longitude" binding:"required"`
	PostalCode    int    `json:"postal_code" binding:"required"`
	IsActive      bool   `json:"is_active"`
}

type PharmacyResponse struct {
	ID            int    `json:"id"`
	PartnerID     int    `json:"partner_id"`
	PartnerName   string `json:"partner_name"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	ProvinceID    int    `json:"province_id"`
	Province      string `json:"province"`
	CityID        int    `json:"city_id"`
	City          string `json:"city"`
	DistrictID    int    `json:"district_id"`
	District      string `json:"district"`
	SubDistrictID int    `json:"sub_district_id"`
	SubDistrict   string `json:"sub_district"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	PostalCode    int    `json:"postal_code"`
	IsActive      bool   `json:"is_active"`
	Logo          string `json:"logo"`
	UpdatedAt     string `json:"updated_at"`
}

type PaginatedPharmaciesResponse struct {
	Pharmacies []PharmacyResponse            `json:"pharmacies"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}

type PharmacyFilterRequest struct {
	Field string `form:"field"`
	Order string `form:"order"`
	Name  string `form:"name"`
	City  string `form:"city"`
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
}

type FileRequest struct {
	File multipart.FileHeader `form:"file"`
	ID   int                  `form:"id"`
}
