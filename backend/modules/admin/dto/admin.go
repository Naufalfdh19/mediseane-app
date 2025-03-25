package dto

type GetusersResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	ProfilePhoto string `json:"profile_photo"`
	Role         string `json:"role"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalUser   int `json:"total_user"`
}

type UserList struct {
	Pagination Pagination         `json:"pagination"`
	Users      []GetusersResponse `json:"user_list"`
}
