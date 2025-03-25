package handler

import (
	"montelukast/modules/delivery/converter"
	"montelukast/modules/delivery/dto"
	"montelukast/modules/delivery/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	deliveryUsecase usecase.DeliveryUsecase
}

func NewDeliveryHandler(deliveryUsecase usecase.DeliveryUsecase) DeliveryHandler {
	return DeliveryHandler{
		deliveryUsecase: deliveryUsecase,
	}
}

func (h *DeliveryHandler) GetOngkirCost(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		err := apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer)
		c.Error(err)
		return
	}
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}
	pharmacyIDStr, exist := c.GetQuery("pharmacy_id")
	if !exist {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrOngkir, apperror.ErrIdEmpty, err))
		return
	}
	pharmacy_id, err := strconv.Atoi(pharmacyIDStr)
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}
	var dtoOngkir []dto.OngkirResponseDTO
	var converter converter.OngkirConverterImpl
	result, err := h.deliveryUsecase.GetAllOngkir(c, userID, pharmacy_id)
	if err != nil {
		c.Error(err)
		return
	}
	for _, data := range result {
		dtoOngkir = append(dtoOngkir, converter.ToDTO(data))
	}
	response := wrapper.ResponseData(dtoOngkir, "get delivery success!", nil)
	c.JSON(http.StatusOK, response)
}
