package converter

import (
	"montelukast/modules/pharmacy/dto"
	"montelukast/modules/pharmacy/entity"
)

type PharmacyConverter interface {
	ToEntity(dto.Pharmacy) entity.Pharmacy
	ToDto(entity.Pharmacy) dto.Pharmacy
}

type PharmacyFilter interface {
	ToEntity(dto.PharmacyFilterRequest) entity.PharmacyFilter
}

type PaginatedPharmaciesConverter interface {
	ToDTO(entity.PaginatedPharmacies) dto.PaginatedPharmaciesResponse
}

type PharmacyFileConverter interface {
	ToEntity(dto.FileRequest) entity.File
}

type PharmacyConverterImpl struct{}
type PharmacyFilterConverterImpl struct{}

func (c *PharmacyFilterConverterImpl) ToEntity(filterDTO dto.PharmacyFilterRequest) (filter entity.PharmacyFilter) {
	return entity.PharmacyFilter{
		Field: filterDTO.Field,
		Order: filterDTO.Order,
		Name:  filterDTO.Name,
		City:  filterDTO.City,
		Limit: filterDTO.Limit,
		Page:  filterDTO.Page,
	}
}

func (c *PharmacyConverterImpl) ToEntity(pharmacyDTO dto.Pharmacy) (pharmacy entity.Pharmacy) {
	return entity.Pharmacy{
		Name:          pharmacyDTO.Name,
		ID:            pharmacyDTO.ID,
		PartnerID:     pharmacyDTO.PartnerID,
		Address:       pharmacyDTO.Address,
		Province:      pharmacyDTO.Province,
		ProvinceID:    pharmacyDTO.ProvinceID,
		District:      pharmacyDTO.District,
		DistrictID:    pharmacyDTO.DistrictID,
		SubDistrict:   pharmacyDTO.SubDistrict,
		SubDistrictID: pharmacyDTO.SubDistrictID,
		City:          pharmacyDTO.City,
		CityID:        pharmacyDTO.CityID,
		Latitude:      pharmacyDTO.Latitude,
		Longitude:     pharmacyDTO.Longitude,
		PostalCode:    pharmacyDTO.PostalCode,
		IsActive:      pharmacyDTO.IsActive,
	}
}

func (c *PharmacyConverterImpl) ToDTO(pharmacy entity.Pharmacy) (pharmacyDTO dto.PharmacyResponse) {
	return dto.PharmacyResponse{
		ID:            pharmacy.ID,
		Name:          pharmacy.Name,
		PartnerID:     pharmacy.PartnerID,
		PartnerName:   pharmacy.PartnerName,
		Address:       pharmacy.Address,
		ProvinceID:    pharmacy.ProvinceID,
		Province:      pharmacy.Province,
		CityID:        pharmacy.CityID,
		City:          pharmacy.City,
		DistrictID:    pharmacy.DistrictID,
		District:      pharmacy.District,
		SubDistrictID: pharmacy.SubDistrictID,
		SubDistrict:   pharmacy.SubDistrict,
		Latitude:      pharmacy.Latitude,
		Longitude:     pharmacy.Longitude,
		PostalCode:    pharmacy.PostalCode,
		IsActive:      pharmacy.IsActive,
		Logo:          pharmacy.Logo,
		UpdatedAt:     pharmacy.UpdatedAt,
	}
}
