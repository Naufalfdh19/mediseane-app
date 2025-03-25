package handler

import (
	"montelukast/modules/pharmacist/converter"
	"montelukast/modules/pharmacist/dto"
	queryparams "montelukast/modules/pharmacist/query_params"
	"montelukast/modules/pharmacist/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PharmacistHandler struct {
	u usecase.PharmacistUsecase
}

func NewPharmacistHandler(u usecase.PharmacistUsecase) PharmacistHandler {
	return PharmacistHandler{
		u: u,
	}
}

func (h *PharmacistHandler) Login(c *gin.Context) {
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

func (h *PharmacistHandler) AddPharmacistHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldAddPharmacist, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	userReq := dto.AddPharmacistRequest{}

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.AddPharmacist(c, converter.AddPharmacistConverter{}.ToEntity(userReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "create pharmacist success!", nil)
	c.JSON(http.StatusCreated, response)
}

func (h *PharmacistHandler) UpdatePharmacistHandler(c *gin.Context) {
	pharmacistIDStr := c.Param("id")
	pharmacistID, err := strconv.Atoi(pharmacistIDStr)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldUpdatePharmacist, apperror.ErrConvertVariableType, nil)
		c.Error(err)
		return
	}

	err = apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldUpdatePharmacist, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	userReq := dto.UpdatePharmacistRequest{}

	err = c.ShouldBindJSON(&userReq)
	if err != nil {
		c.Error(err)
		return
	}

	userEntity := converter.UpdatePharmacistConverter{}.ToEntity(userReq)
	userEntity.ID = pharmacistID
	err = h.u.UpdatePharmacist(c, userEntity)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update pharmacist success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *PharmacistHandler) UpdatePharmacistPhotoHandler(c *gin.Context) {
	pharmacistIDStr := c.Param("id")
	pharmacistID, err := strconv.Atoi(pharmacistIDStr)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldUpdatePharmacistPhoto, apperror.ErrConvertVariableType, nil)
		c.Error(err)
		return
	}
	
	_, fileHeader, err := c.Request.FormFile("file")
	fileName := fileHeader.Filename
	extension := fileName[strings.Index(fileName, ".")+1:]
	if err != nil {
		response := wrapper.ResponseData(nil, "Select a file to upload", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if extension != "jpeg" && extension != "png" && extension != "jpg" {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageType, apperror.ErrUploadImage, err)
		c.Error(err)
		return
	}
	if fileHeader.Size > appconstant.IMAGESIZEMAX {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageSize, apperror.ErrUploadImageSize, err)
		c.Error(err)
		return
	}
	formFile, err := fileHeader.Open()
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageType, apperror.ErrUploadImage, err)
		c.Error(err)
		return
	}

	fileRequest := dto.FileRequest{}
	fileRequest.File = formFile
	var fileConverter converter.PharmacistFileConverter
	url, err := h.u.UpdatePharmacistPhoto(c, fileConverter.ToEntity(fileRequest), pharmacistID)
	if err != nil {
		c.Error(err)
		return
	}
	response := wrapper.ResponseData(url, "success uploaded file", err)
	c.JSON(http.StatusOK, response)
}

func (h *PharmacistHandler) DeletePharmacistHandler(c *gin.Context) {
	pharmacistIDStr := c.Param("id")
	pharmacistID, err := strconv.Atoi(pharmacistIDStr)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldDeletePharmacist, apperror.ErrConvertVariableType, nil)
		c.Error(err)
		return
	}

	err = h.u.DeletePharmacist(c, pharmacistID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "delete pharmacist success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *PharmacistHandler) GetPharmacistsHandler(c *gin.Context) {
	var param queryparams.QueryParamsExistence

	param.Limit, param.IsLimitExists = c.GetQuery("limit")
	param.Page, param.IsPageExists = c.GetQuery("page")
	param.Name, param.IsNameExists = c.GetQuery("name")
	param.Email, param.IsEmailExists = c.GetQuery("email")
	param.SortBy, param.IsSortByExists = c.GetQuery("sort_by")
	param.Order, param.IsOrderExists = c.GetQuery("order")
	param.SipaNumber, param.IsSipaNumberExists = c.GetQuery("sipa_number")
	param.PhoneNumber, param.IsPhoneNumberExists = c.GetQuery("phone_number")
	param.YearOfExperience, param.IsYearOfExperienceExists = c.GetQuery("year_of_experience")

	var queryParams queryparams.QueryParams

	queryParams.Limit = appconstant.LimitInitialvaluePharmacist
	queryParams.Page = appconstant.PageInitialvaluePharmacist
	queryParams.SortBy = appconstant.SortbyInitialvaluePharmacist
	queryParams.Order = appconstant.OrderInitialvaluePharmacist

	h.GetPharmacists(c, queryParams, param)
}

func (h *PharmacistHandler) GetPharmacists(c *gin.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) {
	pharmacists, err := h.u.GetPharmacists(c, queryParams, param)
	if err != nil {
		c.Error(err)
		return
	}
	var pharmacistsDto dto.PharmacistList
	pharmacistsDto.Pagination = converter.PaginationConverter{}.ToDto(pharmacists.Pagination)

	for _, pharmacist := range pharmacists.Pharmacists {
		pharmacistsDto.Pharmacists = append(pharmacistsDto.Pharmacists, converter.GetPharmacistsConverter{}.ToDto(pharmacist))
	}

	response := wrapper.ResponseData(pharmacistsDto, "get pharmacist success!", nil)
	c.JSON(http.StatusOK, response)
}


func (h *PharmacistHandler) GetRandomPassHandler(c *gin.Context) {
	randomPass, err := h.u.GetRandomPass(c)
	if err != nil {
		c.Error(err)
		return
	}
	randomPassDto := converter.GetRandomPassConverter{}.ToDto(*randomPass)

	response := wrapper.ResponseData(randomPassDto, "get random password success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *PharmacistHandler) GetPharmacistDetailHandler(c *gin.Context) {
	rawID := c.Param("id")
	pharmacistID, err := strconv.Atoi(rawID)
	if err != nil {
		err = apperror.NewErrStatusBadRequest(appconstant.FieldErrGetPartner, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	pharmacist, err := h.u.GetPharmacistDetail(c, pharmacistID)
	if err != nil {
		c.Error(err)
		return
	}
	pharmacistDto := converter.GetPharmacistsConverter{}.ToDto(*pharmacist)

	response := wrapper.ResponseData(pharmacistDto, "get pharmacist success!", nil)
	c.JSON(http.StatusOK, response)
}
