package entity

import "mime/multipart"

type Pharmacist struct {
	ID               int
	PharmacyID       *int
	PharmacyName     *string
	Name             string
	SipaNumber       string
	PhoneNumber      string
	YearOfExperience int
	Email            string
	Password         string
	ProfilePhoto     string
	Role             string
}
type Pagination struct {
	CurrentPage     int
	TotalPage       int
	TotalPharmacist int
}

type PharmacistList struct {
	Pagination  Pagination
	Pharmacists []Pharmacist
}

type File struct {
	File multipart.File
}
