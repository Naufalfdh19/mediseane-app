package entity


type User struct {
	ID int
	Name string
	Email string
	ProfilePhoto string
	Role string
}

type Pagination struct {
	CurrentPage int
	TotalPage   int
	TotalUser   int
}

type UserList struct {
	Pagination  Pagination
	Users []User
}