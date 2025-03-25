package handler

import (
	"montelukast/modules/address/converter"
	"montelukast/modules/address/dto"
	"montelukast/modules/address/entity"
	"montelukast/modules/address/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	u usecase.AddressUsecase
}

func NewAddressHandler(u usecase.AddressUsecase) AddressHandler {
	return AddressHandler{
		u: u,
	}
}

func (h *AddressHandler) GetProvincesHandler(c *gin.Context) {
	provinces, err := h.u.GetProvinces(c)
	if err != nil {
		c.Error(err)
		return
	}

	provincesDto := make([]dto.ProvinceResponse, 0)

	for _, province := range provinces {
		provincesDto = append(provincesDto, converter.GetProvincesConverter{}.ToDto(province))
	}

	response := wrapper.ResponseData(provincesDto, "get provinces success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) GetCitiesHandler(c *gin.Context) {
	provinceIDStr, ok := c.GetQuery("province")
	if !ok {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCities, apperror.ErrQueryParams, apperror.ErrQueryParams))
		return
	}
	provinceID, err := strconv.Atoi(provinceIDStr)
	if err != nil {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCities, apperror.ErrQueryParams, err))
		return
	}

	cities, err := h.u.GetCities(c, provinceID)
	if err != nil {
		c.Error(err)
		return
	}

	citiesDto := make([]dto.CityAndDistrictResponse, 0)

	for _, city := range cities {
		citiesDto = append(citiesDto, converter.GetCitiesAndDistrictsConverter{}.ToDto(city))
	}

	response := wrapper.ResponseData(citiesDto, "get cities success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) GetDistrictsHandler(c *gin.Context) {
	cityIDStr, ok := c.GetQuery("city")
	if !ok {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetDistricts, apperror.ErrQueryParams, apperror.ErrQueryParams))
		return
	}
	cityID, err := strconv.Atoi(cityIDStr)
	if err != nil {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetDistricts, apperror.ErrQueryParams, err))
		return
	}

	districts, err := h.u.GetDistricts(c, cityID)
	if err != nil {
		c.Error(err)
		return
	}

	districtsDto := make([]dto.CityAndDistrictResponse, 0)

	for _, district := range districts {
		districtsDto = append(districtsDto, converter.GetCitiesAndDistrictsConverter{}.ToDto(district))
	}

	response := wrapper.ResponseData(districtsDto, "get districts success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) GetSubDistrictsHandler(c *gin.Context) {
	districtIDstr, ok := c.GetQuery("district")
	if !ok {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetSubDistricts, apperror.ErrQueryParams, apperror.ErrQueryParams))
		return
	}
	districtID, err := strconv.Atoi(districtIDstr)
	if err != nil {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetSubDistricts, apperror.ErrQueryParams, err))
		return
	}

	subDistricts, err := h.u.GetSubDistricts(c, districtID)
	if err != nil {
		c.Error(err)
		return
	}

	subDistrictsDto := make([]dto.SubDistrictResponse, 0)

	for _, subDistrict := range subDistricts {
		subDistrictsDto = append(subDistrictsDto, converter.GetSubDistrictsConverter{}.ToDto(subDistrict))
	}

	response := wrapper.ResponseData(subDistrictsDto, "get sub-districts success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) AddUserAddressHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrAddAddress, apperror.ErrInvalidJSON, err))
		return
	}

	addUserAddressReq := dto.AddUserAddressRequest{}
	err = c.ShouldBindJSON(&addUserAddressReq)
	if err != nil {
		c.Error(err)
		return
	}

	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	err = h.u.AddUserAddress(c, converter.AddUserAddressConverter{}.ToEntity(addUserAddressReq, userID))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "add address success!", nil)
	c.JSON(http.StatusCreated, response)
}

func (h *AddressHandler) GetCurrentLocationHandler(c *gin.Context) {
	longitude, ok := c.GetQuery("long")
	if !ok {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCurrentLocation, apperror.ErrQueryParams, apperror.ErrQueryParams))
		return
	}
	latitude, ok := c.GetQuery("lat")
	if !ok {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCurrentLocation, apperror.ErrQueryParams, apperror.ErrQueryParams))
		return
	}

	address, err := h.u.GetCurrentLocation(c, longitude, latitude)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(converter.GetCurrentLocationConverter{}.ToDto(*address), "get current location success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) GetUserAddressesHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	filterDTO := dto.GetUserAddressFilter{}
	if err := c.ShouldBindQuery(&filterDTO); err != nil {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetAddresses, apperror.ErrConvertVariableType, err))
		return
	}

	addresses, err := h.u.GetUserAddresses(c, userID, entity.AddressFilter(filterDTO))
	if err != nil {
		c.Error(err)
		return
	}

	addressesDto := make([]dto.GetUserAddressesResponse, 0)

	for _, address := range addresses {
		addressesDto = append(addressesDto, converter.GetUserAddressesConverter{}.ToDto(address))
	}

	response := wrapper.ResponseData(addressesDto, "get user addresses success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) UpdateUserAddressHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	err = apperror.JsonValidator(c)
	if err != nil {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrAddAddress, apperror.ErrInvalidJSON, err))
		return
	}
	updateUserAddressReq := dto.UpdateUserAddressRequest{}
	err = c.ShouldBindJSON(&updateUserAddressReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.UpdateUserAddress(c, converter.UpdateUserAddressConverter{}.ToDto(updateUserAddressReq, userID))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update address success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) GetUserAddressHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	addressIDStr := c.Param("id")
	addressID, err := strconv.Atoi(addressIDStr)
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	address, err := h.u.GetUserAddress(c, userID, addressID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(converter.GetUserAddressConverter{}.ToDto(*address), "get user address success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) DeleteUserAddressHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	addressIDStr := c.Param("id")
	addressID, err := strconv.Atoi(addressIDStr)
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	err = h.u.DeleteUserAddress(c, userID, addressID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "delete user address success!", nil)
	c.JSON(http.StatusOK, response)
}
