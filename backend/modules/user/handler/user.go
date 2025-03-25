package handler

import (
	"montelukast/modules/user/converter"
	"montelukast/modules/user/dto"
	"montelukast/modules/user/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	u usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) UserHandler {
	return UserHandler{
		u: u,
	}
}

func (h *UserHandler) RegisterHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrRegister, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	userReq := dto.UserRegisterRequest{}

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.Register(c, converter.UserRegisterConverter{}.ToEntity(userReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "register success!", nil)
	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) LoginHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	userReq := dto.UserLoginRequest{}

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := h.u.Login(c, converter.UserLoginConverter{}.ToEntity(userReq))
	if err != nil {
		c.Error(err)
		return
	}
	resDto := dto.UserLoginResponse{
		AccessToken: res,
	}

	response := wrapper.ResponseData(resDto, "login success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) ForgetPasswordHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrForgetPassword, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	userReq := dto.UserForgetPassRequest{}

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.ForgetPassword(c, converter.ForgetPasswordConverter{}.ToEntity(userReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "reset password link, successfully sent!", nil)
	c.JSON(http.StatusCreated, response)
}

func (h UserHandler) CheckResetPassTokenHandler(c *gin.Context) {
	token, tokenExists := c.GetQuery("token")
	if !tokenExists {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, nil)
		c.Error(err)
		return
	}

	err := h.u.CheckResetPassToken(c, token)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "reset password token exists!", nil)
	c.JSON(http.StatusOK, response)
}

func (h UserHandler) ResetPasswordHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrResetPassword, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}

	userReq := dto.UserResetPasswordRequest{}

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.ResetPassword(c, converter.ResetPasswordConverter{}.ToEntity(userReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "reset password success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h UserHandler) VerifyEmailHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrVerifyEmail, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}

	userReq := dto.UserVerifyEmailRequest{}

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.VerifyEmail(c, userReq.Token)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "verify email success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h UserHandler) GetProfileHandler(c *gin.Context) {
	rawUserID, isExists := c.Get("user_id")
	if !isExists {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, apperror.ErrTokenInvalid)
		c.Error(err)
		return
	}
	userID, err := strconv.Atoi(rawUserID.(string))
	if err != nil {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, err)
		c.Error(err)
		return
	}

	user, err := h.u.GetProfile(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(converter.UserGetProfileConverter{}.ToDTO(*user), "get user profile success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateNameHandler(c *gin.Context) {
	rawUserID, isExists := c.Get("user_id")
	if !isExists {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, apperror.ErrTokenInvalid)
		c.Error(err)
		return
	}
	userID, err := strconv.Atoi(rawUserID.(string))
	if err != nil {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, err)
		c.Error(err)
		return
	}

	err = apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateUser, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	userReq := dto.UserUpdateNameRequest{}
	userReq.ID = userID

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.UpdateName(c, converter.UserUpdateNameConverter{}.ToEntity(userReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update username success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) SendEmailFromProfileHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrEmail, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}

	userReq := dto.UserForgetPassRequest{}

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}
	err = h.u.ResendVerificationEmail(c, userReq.Email)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "send email success!", nil)
	c.JSON(http.StatusOK, response)
}

