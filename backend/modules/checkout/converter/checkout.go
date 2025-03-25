package converter

import (
	"montelukast/modules/checkout/dto"
	"montelukast/modules/checkout/entity"
)

type CheckoutConverter interface {
	ToEntity(dto.CheckoutData) entity.CheckoutData
}

type ListDeliveryConverter interface {
	ToEntity(dto.DeliveryData) entity.DeliveryData
}

type CheckoutConverterImpl struct{}
type ListDeliveryConverterImpl struct{}

func (c CheckoutConverterImpl) ToEntity(dto dto.CheckoutData) entity.CheckoutData {
	var converter ListDeliveryConverterImpl
	var entityList []entity.DeliveryData
	for _, data := range dto.ListDeliveryData {
		entityList = append(entityList, converter.ToEntity(data))
	}
	return entity.CheckoutData{
		IDCart:           dto.IDCart,
		ListDeliveryData: entityList,
	}
}

func (c ListDeliveryConverterImpl) ToEntity(dto dto.DeliveryData) entity.DeliveryData {
	return entity.DeliveryData{
		PharmacyID: dto.PharmacyID,
		DeliveryID: dto.DeliveryID,
	}
}
