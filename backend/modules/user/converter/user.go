package converter

import (
	"montelukast/modules/user/dto"
	"montelukast/modules/user/entity"
)

type UserRegisterConverter struct{}

func (c UserRegisterConverter) ToEntity(userReq dto.UserRegisterRequest) entity.User {
	return entity.User{
		Name:     userReq.Username,
		Email:    userReq.Email,
		Password: userReq.Password,
	}
}

type UserLoginConverter struct{}

func (c UserLoginConverter) ToEntity(userReq dto.UserLoginRequest) entity.User {
	return entity.User{
		Email:    userReq.Email,
		Password: userReq.Password,
	}
}

type ForgetPasswordConverter struct{}

func (c ForgetPasswordConverter) ToEntity(userReq dto.UserForgetPassRequest) entity.User {
	return entity.User{
		Email: userReq.Email,
	}
}

type ResetPasswordConverter struct{}

func (r ResetPasswordConverter) ToEntity(userReq dto.UserResetPasswordRequest) entity.ResetPassword {
	return entity.ResetPassword{
		Token:       userReq.Token,
		NewPassword: userReq.NewPassword,
	}
}

type VerifyEmailConverter struct{}

func (v VerifyEmailConverter) ToEntity(userReq dto.UserVerifyEmailRequest) entity.ResetPassword {
	return entity.ResetPassword{
		Token: userReq.Token,
	}
}

type UserGetProfileConverter struct{}

func (p UserGetProfileConverter) ToDTO(user entity.User) dto.UserGetProfileResponse {
	return dto.UserGetProfileResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		ProfilePhoto: user.ProfilePhoto,
		Role:         user.Role,
		IsVerified:   user.IsVerified,
	}
}

type UserUpdateNameConverter struct{}

func (c UserUpdateNameConverter) ToEntity(userReq dto.UserUpdateNameRequest) entity.User {
	return entity.User{
		ID:   userReq.ID,
		Name: userReq.Username,
	}
}
