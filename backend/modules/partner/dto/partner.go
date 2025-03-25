package dto

type AddPartnerRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	YearFounded string `json:"year_founded" binding:"required,gte=0"`
	ActiveDays  string `json:"active_days" binding:"required"`
	StartHour   string `json:"start_hour" binding:"required"`
	EndHour     string `json:"end_hour" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

type UpdatePartnerRequest struct {
	ActiveDays string `json:"active_days" binding:"required"`
	StartHour  string `json:"start_hour" binding:"required"`
	EndHour    string `json:"end_hour" binding:"required"`
	IsActive   bool   `json:"is_active"`
}

type DeletePartnerRequest struct {
	PartnerID int `json:"partner_id" binding:"required"`
}

type GetPartnersResponse struct {
	ID          int    `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	YearFounded string `json:"year_founded" binding:"required"`
	ActiveDays  string `json:"active_days" binding:"required"`
	StartHour   string `json:"start_hour" binding:"required"`
	EndHour     string `json:"end_hour" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

type Pagination struct {
	CurrentPage     int `json:"current_page"`
	TotalPage       int `json:"total_page"`
	TotalPharmacist int `json:"total_item"`
}

type PartnerList struct {
	Pagination Pagination            `json:"pagination"`
	Partners   []GetPartnersResponse `json:"partner_list"`
}
