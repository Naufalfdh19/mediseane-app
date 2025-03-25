package entity

type Partner struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	YearFounded string `json:"year_founded"`
	ActiveDays  string `json:"active_days"`
	StartHour   string `json:"start_hour"`
	EndHour     string `json:"end_hour"`
	IsActive    bool   `json:"is_active"`
}

type Pagination struct {
	CurrentPage     int
	TotalPage       int
	TotalPharmacist int
}

type PartnerList struct {
	Pagination Pagination
	Partners   []Partner
}
