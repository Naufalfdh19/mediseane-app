package handler

import (
	adminConverter "montelukast/modules/admin/converter"
	adminDto "montelukast/modules/admin/dto"
	queryparams "montelukast/modules/admin/query_params"
	"montelukast/modules/admin/usecase"
	"montelukast/modules/user/converter"
	"montelukast/modules/user/dto"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	u usecase.AdminUsecase
}

func NewAdminHandler(u usecase.AdminUsecase) AdminHandler {
	return AdminHandler{
		u: u,
	}
}

func (h *AdminHandler) Login(c *gin.Context) {
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

func (h *AdminHandler) GetUsersHandler(c *gin.Context) {
	var param queryparams.QueryParamsExistence

	param.Name, param.IsNameExists = c.GetQuery("name")
	param.Email, param.IsEmailExists = c.GetQuery("email")
	param.Role, param.IsRoleExists = c.GetQuery("role")
	param.Limit, param.IsLimitExists = c.GetQuery("limit")
	param.Page, param.IsPageExists = c.GetQuery("page")
	param.SortBy, param.IsSortByExsts = c.GetQuery("sort_by")
	param.Order, param.IsOrderExists = c.GetQuery("order")

	var queryParams queryparams.QueryParams

	queryParams.Limit = appconstant.LimitInitialvaluePharmacist
	queryParams.Page = appconstant.PageInitialvaluePharmacist
	queryParams.SortBy = appconstant.SortbyInitialvaluePharmacist
	queryParams.Order = appconstant.OrderInitialvaluePharmacist

	h.GetUsers(c, queryParams, param)
}

func (h *AdminHandler) GetUsers(c *gin.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) {
	users, err := h.u.GetUsers(c, queryParams, param)
	if err != nil {
		c.Error(err)
		return
	}
	var usersDto adminDto.UserList
	usersDto.Pagination = adminConverter.PaginationConverter{}.ToDto(users.Pagination)

	for _, user := range users.Users {
		usersDto.Users = append(usersDto.Users, adminConverter.GetusersConverter{}.ToDto(user))
	}

	response := wrapper.ResponseData(usersDto, "get users success!", nil)
	c.JSON(http.StatusOK, response)
}
