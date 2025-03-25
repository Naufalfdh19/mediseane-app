package dto

type ProvinceResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CityAndDistrictResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Longitude *string `json:"longitude"`
	Latitude  *string `json:"latitude"`
}

type SubDistrictResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PostalCodes string `json:"postal_codes"`
}

type AddUserAddressRequest struct {
	Name          string `json:"name" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	Address       string `json:"address" binding:"required"`
	ProvinceID    int    `json:"province_id" binding:"required"`
	Province      string `json:"province" binding:"required"`
	CityID        int    `json:"city_id" binding:"required"`
	City          string `json:"city" binding:"required"`
	DistrictID    int    `json:"district_id" binding:"required"`
	District      string `json:"district" binding:"required"`
	SubDistrictID int    `json:"sub_district_id" binding:"required"`
	SubDistrict   string `json:"sub_district" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	Longitude     string `json:"longitude" binding:"required"`
	Latitude      string `json:"latitude" binding:"required"`
}

type GetCurrentLocationResponse struct {
	ProvinceId int `json:"province_id"`
	CityId     int `json:"city_id"`
}

type GetUserAddressesResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	SubDistrict string `json:"sub_district"`
	PostalCode  string `json:"postal_code"`
	IsActive    bool   `json:"is_active"`
}

type GetUserAddressFilter struct {
	IsActive string `form:"active"`
}

type UpdateUserAddressRequest struct {
	ID            int    `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	Address       string `json:"address" binding:"required"`
	ProvinceID    int    `json:"province_id" binding:"required"`
	Province      string `json:"province" binding:"required"`
	CityID        int    `json:"city_id" binding:"required"`
	City          string `json:"city" binding:"required"`
	DistrictID    int    `json:"district_id" binding:"required"`
	District      string `json:"district" binding:"required"`
	SubDistrictID int    `json:"sub_district_id" binding:"required"`
	SubDistrict   string `json:"sub_district" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	Longitude     string `json:"longitude" binding:"required"`
	Latitude      string `json:"latitude" binding:"required"`
	IsActive      *bool  `json:"is_active" binding:"required"`
}

type GetUserAddressResponse struct {
	Name          string `json:"name"`
	PhoneNumber   string `json:"phone_number"`
	Address       string `json:"address"`
	ProvinceID    int    `json:"province_id"`
	CityID        int    `json:"city_id"`
	DistrictID    int    `json:"district_id"`
	SubDistrictID int    `json:"sub_district_id"`
	PostalCode    string `json:"postal_code"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	IsActive      bool   `json:"is_active"`
}
