package converter

import (
	"montelukast/modules/delivery/dto"
	"montelukast/modules/delivery/entity"
)

type OngkirConverter interface {
	ToDTO(entity.OngkirData) dto.OngkirCostResponse
}

type OngkirConverterImpl struct{}

func (c OngkirConverterImpl) ToDTO(entity entity.OngkirData) dto.OngkirResponseDTO {
	return dto.OngkirResponseDTO{
		Id:   entity.Id,
		Name: entity.Name,
		Cost: entity.Cost,
		Etd:  entity.Etd,
	}
}
