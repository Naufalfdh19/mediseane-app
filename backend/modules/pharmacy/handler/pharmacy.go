package handler

import (
	"montelukast/modules/pharmacy/converter"
	"montelukast/modules/pharmacy/dto"
	"montelukast/modules/pharmacy/entity"
	"montelukast/modules/pharmacy/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/pagination"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PharmacyHandler struct {
	pharmacyUsecase usecase.PharmacyUsecase
}

func NewPharmacyHandler(pharmacyUsecase usecase.PharmacyUsecase) PharmacyHandler {
	return PharmacyHandler{
		pharmacyUsecase: pharmacyUsecase,
	}
}

func (h PharmacyHandler) GetPharmacyByID(c *gin.Context) {
	id := c.Param("id")
	pharmacyId, err := strconv.Atoi(id)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrPharmacy, apperror.ErrIdType, err)
		c.Error(err)
		return
	}
	var pharmacyConverter converter.PharmacyConverterImpl
	result, err := h.pharmacyUsecase.GetPharmacyByID(c, pharmacyId)
	if err != nil {
		c.Error(err)
		return
	}
	response := wrapper.ResponseData(pharmacyConverter.ToDTO(result), "get pharmacy success", nil)
	c.JSON(http.StatusOK, response)
}

func (h PharmacyHandler) GetAllPharmaciesHandler(c *gin.Context) {
	filterDTO := dto.PharmacyFilterRequest{}
	if err := c.ShouldBindQuery(&filterDTO); err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetPharmacies, apperror.ErrConvertVariableType, nil)
		c.Error(err)
		return
	}
	var filterConverter converter.PharmacyFilterConverterImpl
	pharmacies, err := h.pharmacyUsecase.GetAllPharmacies(c, filterConverter.ToEntity(filterDTO))
	if err != nil {
		c.Error(err)
		return
	}
	results := []dto.PharmacyResponse{}

	var pharmacyConverter converter.PharmacyConverterImpl
	for _, pharmacy := range pharmacies.Pharmacies {
		results = append(results, pharmacyConverter.ToDTO(pharmacy))
	}
	var finalResult dto.PaginatedPharmaciesResponse
	finalResult.Pharmacies = results
	var paginationInfo pagination.PaginationConverter
	finalResult.Pagination = paginationInfo.ToDto(pharmacies.Pagination)

	response := wrapper.ResponseData(finalResult, "get pharmacies success", nil)
	c.JSON(http.StatusOK, response)
}

func (h PharmacyHandler) AddPharmacyHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrJSON, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	var pharmacyConverter converter.PharmacyConverterImpl
	pharmacyData := dto.Pharmacy{}
	err = c.ShouldBind(&pharmacyData)
	if err != nil {
		c.Error(err)
		return
	}
	err = h.pharmacyUsecase.AddPharmacy(c, pharmacyConverter.ToEntity(pharmacyData))
	if err != nil {
		c.Error(err)
		return
	}
	response := wrapper.ResponseData(nil, "add pharmacy success", nil)
	c.JSON(http.StatusCreated, response)
}

func (h PharmacyHandler) UpdatePharmacyHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrJSON, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	var pharmacyConverter converter.PharmacyConverterImpl
	pharmacyData := dto.Pharmacy{}
	err = c.ShouldBind(&pharmacyData)
	if err != nil {
		c.Error(err)
		return
	}
	err = h.pharmacyUsecase.UpdatePharmacy(c, pharmacyConverter.ToEntity(pharmacyData))
	if err != nil {
		c.Error(err)
		return
	}
	response := wrapper.ResponseData(nil, "update pharmacy success", nil)
	c.JSON(http.StatusOK, response)
}

func (h PharmacyHandler) DeletePharmacyHandler(c *gin.Context) {
	id := c.Param("id")
	pharmacyId, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		return
	}
	err = h.pharmacyUsecase.DeletePharmacy(c, pharmacyId)
	if err != nil {
		c.Error(err)
		return
	}
	response := wrapper.ResponseData(nil, "delete pharmacy success", nil)
	c.JSON(http.StatusOK, response)
}

func (h PharmacyHandler) AddLogoHandler(c *gin.Context) {
	fileRequest := dto.FileRequest{}
	err := c.ShouldBind(&fileRequest)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetPharmacies, apperror.ErrConvertVariableType, nil)
		c.Error(err)
	}
	fileName := fileRequest.File.Filename
	extension := fileName[strings.Index(fileName, ".")+1:]
	if extension != "jpeg" && extension != "png" && extension != "jpg" {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageType, apperror.ErrUploadImage, err)
		c.Error(err)
		return
	}
	if fileRequest.File.Size > appconstant.IMAGESIZEMAX {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageSize, apperror.ErrUploadImageSize, err)
		c.Error(err)
		return
	}
	formFile, err := fileRequest.File.Open()
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageType, apperror.ErrUploadImage, err)
		c.Error(err)
		return
	}

	var entityFile entity.File
	entityFile.File = formFile
	entityFile.ID = fileRequest.ID
	url, err := h.pharmacyUsecase.AddLogo(c, entityFile)
	if err != nil {
		c.Error(err)
		return
	}
	response := wrapper.ResponseData(url, "success uploaded file", err)
	c.JSON(http.StatusOK, response)
}
