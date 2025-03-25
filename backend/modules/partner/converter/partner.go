package converter

import (
	"montelukast/modules/partner/dto"
	"montelukast/modules/partner/entity"
)


type AddPartnerConverter struct {}

func (c AddPartnerConverter) ToEntity(partnerReq dto.AddPartnerRequest) entity.Partner {
	return entity.Partner{
		Name: partnerReq.Name,
		YearFounded: partnerReq.YearFounded,
		ActiveDays: partnerReq.ActiveDays,
		StartHour: partnerReq.StartHour,
		EndHour: partnerReq.EndHour,
		IsActive: partnerReq.IsActive,
	}
}

type UpdatePartnerConverter struct {}

func (c UpdatePartnerConverter) ToEntity(partnerReq dto.UpdatePartnerRequest) entity.Partner {
	return entity.Partner{
		ActiveDays: partnerReq.ActiveDays,
		StartHour: partnerReq.StartHour,
		EndHour: partnerReq.EndHour,
		IsActive: partnerReq.IsActive,

	}
}

type PaginationConverter struct{}

func (c PaginationConverter) ToDto(pagination entity.Pagination) dto.Pagination {
	return dto.Pagination{
		CurrentPage: pagination.CurrentPage,
		TotalPage: pagination.TotalPage,
		TotalPharmacist: pagination.TotalPharmacist,
	}
}

type GetPartnersConverter struct{}

func (c GetPartnersConverter) ToDto(partner entity.Partner) dto.GetPartnersResponse {
	return dto.GetPartnersResponse{
		ID: partner.ID,
		Name: partner.Name,
		YearFounded: partner.YearFounded,
		ActiveDays: partner.ActiveDays,
		StartHour: partner.StartHour,
		EndHour: partner.EndHour,
		IsActive: partner.IsActive,
	}
}


type DeletePartnerConverter struct {}

func (c DeletePartnerConverter) ToEntity(partnerReq dto.DeletePartnerRequest) entity.Partner {
	return entity.Partner{
		ID: partnerReq.PartnerID,
	}
}