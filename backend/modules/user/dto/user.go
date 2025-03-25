package dto

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,gte=5,lte=12"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
}

type UserForgetPassRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UserResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UserVerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type UserGetProfileResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	ProfilePhoto string `json:"profile_photo"`
	Role         string `json:"role"`
	IsVerified   bool   `json:"is_verified"`
}

type UserUpdateNameRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username" binding:"required,gte=5,lte=12"`
}
