package entity

type Location struct {
	ID          int
	Name        string
	Longitude   *string
	Latitude    *string
	PostalCodes string
}

type UserAddress struct {
	ID            int
	UserID        int
	Name          string
	PhoneNumber   string
	Address       string
	ProvinceID    int
	Province      string
	CityID        int
	City          string
	DistrictID    int
	District      string
	SubDistrictID int
	SubDistrict   string
	PostalCode    string
	Longitude     string
	Latitude      string
	IsActive      bool
}
