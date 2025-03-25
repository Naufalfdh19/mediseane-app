package usecase

import (
	"context"
	"math"
	adminEntity "montelukast/modules/admin/entity"
	queryparams "montelukast/modules/admin/query_params"
	"montelukast/modules/admin/repository"
	"montelukast/modules/user/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	jwttoken "montelukast/pkg/jwt_token"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase interface {
	Login(c context.Context, admin entity.User) (token string, err error)
	GetUsers(c context.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) (*adminEntity.UserList, error)
}

type adminUsecaseImpl struct {
	r repository.AdminRepo
}

func NewAdminUsecase(r repository.AdminRepo) adminUsecaseImpl {
	return adminUsecaseImpl{
		r: r,
	}
}

func (u adminUsecaseImpl) Login(c context.Context, user entity.User) (token string, err error) {
	jwt_token := jwttoken.NewJWT()
	exist, err := u.r.IsUserExistByEmail(c, user.Email)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidEmailOrPassword, err)
	}
	admin, err := u.r.GetUserByEmail(c, user.Email)
	if err != nil {
		return "", err
	}
	if admin.Role != appconstant.ROLE_ADMIN {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidEmailOrPassword, err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(user.Password)); err != nil {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidEmailOrPassword, err)
	}
	userIdStr := strconv.Itoa(int(admin.ID))
	token, err = jwt_token.GenerateJwtTokenForAuth(appconstant.JwtTokenAuthType, userIdStr, admin.Role)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return token, nil
}

func (u adminUsecaseImpl) GetUsers(c context.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) (*adminEntity.UserList, error) {
	err := checkQueryParams(&queryParams, param)
	if err != nil {
		return nil, err
	}

	queryParams = queryparams.DefaultQuery(queryParams)

	totalUser, err := u.r.GetTotalUser(c, queryParams)
	if err != nil {
		return nil, err
	}

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalPage := int(math.Ceil(float64(totalUser) / float64(queryParams.Limit)))
	if totalUser <= 0 {
		totalPage = 1
	}

	queryParams.Page = apperror.CheckCurrentPage(queryParams.Page, totalPage, queryParams.Limit)

	users, err := u.r.GetUsers(c, queryParams, totalUser)
	if err != nil {
		return nil, err
	}

	pagination := adminEntity.Pagination{
		CurrentPage: queryParams.Page,
		TotalPage:   totalPage,
		TotalUser:   totalUser,
	}

	usersList := adminEntity.UserList{
		Pagination: pagination,
		Users:      users,
	}

	return &usersList, nil
}

func checkQueryParams(queryParams *queryparams.QueryParams, param queryparams.QueryParamsExistence) error {
	var convertError error
	if param.IsLimitExists {
		limitInt, err := strconv.Atoi(param.Limit)
		if err != nil {
			convertError = err
		}
		queryParams.Limit = limitInt
	}
	if param.IsPageExists {
		pageInt, err := strconv.Atoi(param.Page)
		if err != nil {
			convertError = err
		}
		queryParams.Page = pageInt
	}
	if convertError != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetUser, apperror.ErrConvertVariableType, nil)
		return err
	}

	if param.IsNameExists {
		queryParams.Name = "%" + param.Name + "%"
	}

	if param.IsEmailExists {
		queryParams.Email = "%" + param.Email + "%"
	}

	if param.IsRoleExists {
		queryParams.Role = "%" + param.Role + "%"
	}

	if param.IsSortByExsts {
		queryParams.SortBy = param.SortBy
	}
	if param.IsOrderExists {
		queryParams.Order = param.Order
	}

	return nil
}


