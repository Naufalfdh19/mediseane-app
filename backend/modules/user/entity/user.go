package entity

type User struct {
	ID int
	Name string
	Email string
	Password string
	ProfilePhoto string
	Role string
	IsVerified bool
}

type ResetPassword struct {
	Token string
	NewPassword string
}