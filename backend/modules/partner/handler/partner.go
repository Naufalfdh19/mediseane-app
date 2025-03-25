package handler

import (
	"montelukast/modules/partner/converter"
	"montelukast/modules/partner/dto"
	queryparams "montelukast/modules/partner/query_params"
	"montelukast/modules/partner/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PartnerHandler struct {
	u usecase.PartnerUsecase
}

func NewPartnerHandler(u usecase.PartnerUsecase) PartnerHandler {
	return PartnerHandler{
		u: u,
	}
}

func (h PartnerHandler) AddPartnerHandler(c *gin.Context) {
	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPartner, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	partnerReq := dto.AddPartnerRequest{}

	err = c.ShouldBindJSON(&partnerReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.AddPartner(c, converter.AddPartnerConverter{}.ToEntity(partnerReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "create partner success!", nil)
	c.JSON(http.StatusCreated, response)
}

func (h PartnerHandler) UpdatePartnerHandler(c *gin.Context) {
	partnerIDStr := c.Param("id")
	partnerID, err := strconv.Atoi(partnerIDStr)
	if err != nil {
		err = apperror.NewErrStatusBadRequest(appconstant.FieldErrDeletePartner, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	err = apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdatePartner, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	partnerReq := dto.UpdatePartnerRequest{}

	err = c.ShouldBindJSON(&partnerReq)
	if err != nil {
		c.Error(err)
		return
	}

	partnerEnt := converter.UpdatePartnerConverter{}.ToEntity(partnerReq)
	partnerEnt.ID = partnerID
	err = h.u.UpdatePartner(c, partnerEnt)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update partner success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h PartnerHandler) DeletePartnerHandler(c *gin.Context) {
	partnerIDStr := c.Param("id")
	partnerID, err := strconv.Atoi(partnerIDStr)
	if err != nil {
		err = apperror.NewErrStatusBadRequest(appconstant.FieldErrDeletePartner, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	err = h.u.DeletePartner(c, partnerID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "delete partner success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *PartnerHandler) GetPartnersHandler(c *gin.Context) {
	var param queryparams.QueryParamsExistence

	param.Limit, param.IsLimitExists = c.GetQuery("limit")
	param.Page, param.IsPageExists = c.GetQuery("page")
	param.Name, param.IsNameExists = c.GetQuery("name")
	param.YearFounded, param.IsYearFoundedExists = c.GetQuery("year_founded")
	param.ActiveDays, param.IsActiveDaysExists = c.GetQuery("active_days")
	param.StartHour, param.IsStartHourExists = c.GetQuery("start_hour")
	param.EndHour, param.IsEndHourExists = c.GetQuery("end_hour")
	param.IsActive, param.IsActiveExists = c.GetQuery("is_active")
	param.SortBy, param.IsSortByExsts = c.GetQuery("sort_by")
	param.Order, param.IsOrderExists = c.GetQuery("order")

	var queryParams queryparams.QueryParams

	queryParams.Limit = appconstant.LimitInitialvaluePharmacist
	queryParams.Page = appconstant.PageInitialvaluePharmacist
	queryParams.SortBy = appconstant.SortbyInitialvaluePharmacist
	queryParams.Order = appconstant.OrderInitialvaluePharmacist

	h.GetPartners(c, queryParams, param)
}

func (h *PartnerHandler) GetPartners(c *gin.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) {
	partners, err := h.u.GetPartners(c, queryParams, param)
	if err != nil {
		c.Error(err)
		return
	}
	var partnersDto dto.PartnerList
	partnersDto.Pagination = converter.PaginationConverter{}.ToDto(partners.Pagination)

	for _, pharmacist := range partners.Partners {
		partnersDto.Partners = append(partnersDto.Partners, converter.GetPartnersConverter{}.ToDto(pharmacist))
	}

	response := wrapper.ResponseData(partnersDto, "get partners success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *PartnerHandler) GetPartnerHandler(c *gin.Context) {
	partnerIDStr := c.Param("id")
	partnerID, err := strconv.Atoi(partnerIDStr)
	if err != nil {
		err = apperror.NewErrStatusBadRequest(appconstant.FieldErrGetPartner, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	partner, err := h.u.GetPartner(c, partnerID)
	if err != nil {
		c.Error(err)
		return
	}
	partnerDto := converter.GetPartnersConverter{}.ToDto(*partner)

	response := wrapper.ResponseData(partnerDto, "get partner success!", nil)
	c.JSON(http.StatusOK, response)
}
