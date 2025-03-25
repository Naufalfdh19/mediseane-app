package dto

import "mime/multipart"

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
}

type AddPharmacistRequest struct {
	Name             string `json:"name" binding:"required,gte=5,lte=12"`
	SipaNumber       string `json:"sipa_number" binding:"required,min=8,max=19"`
	PhoneNumber      string `json:"phone_number" binding:"required,min=7,max=15"`
	YearOfExperience int    `json:"year_of_experience" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required"`
}

type UpdatePharmacistRequest struct {
	PharmacyID       *int   `json:"pharmacy_id"`
	PhoneNumber      string `json:"phone_number" binding:"min=7,max=15"`
	YearOfExperience int    `json:"year_of_experience"`
}

type UpdatePharmacistPhotoRequest struct {
	ProfilePhoto string `form:"profile_photo" binding:"required"`
}

type GetPharmacistResponse struct {
	ID               int     `json:"id"`
	PharmacyID       *int    `json:"pharmacy_id,omitempty"`
	PharmacyName     *string `json:"pharmacy_name,omitempty"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	SipaNumber       string  `json:"sipa_number"`
	PhoneNumber      string  `json:"phone_number"`
	YearOfExperience int     `json:"year_of_experience"`
}

type Pagination struct {
	CurrentPage     int `json:"current_page"`
	TotalPage       int `json:"total_page"`
	TotalPharmacist int `json:"total_item"`
}

type PharmacistList struct {
	Pagination  Pagination              `json:"pagination"`
	Pharmacists []GetPharmacistResponse `json:"pharmacist_list"`
}

type FileRequest struct {
	File multipart.File `json:"file,omitempty"`
}

type GetRandomPassResponse struct {
	Password string `json:"password"`
}
