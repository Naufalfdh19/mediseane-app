package converter

import (
	"montelukast/modules/pharmacist/dto"
	"montelukast/modules/pharmacist/entity"
)

type UserLoginConverter struct{}

func (c UserLoginConverter) ToEntity(userReq dto.UserLoginRequest) entity.Pharmacist {
	return entity.Pharmacist{
		Email:    userReq.Email,
		Password: userReq.Password,
	}
}

type AddPharmacistConverter struct{}

func (c AddPharmacistConverter) ToEntity(pharmacistReq dto.AddPharmacistRequest) entity.Pharmacist {
	return entity.Pharmacist{
		Name:             pharmacistReq.Name,
		SipaNumber:       pharmacistReq.SipaNumber,
		PhoneNumber:      pharmacistReq.PhoneNumber,
		YearOfExperience: pharmacistReq.YearOfExperience,
		Email:            pharmacistReq.Email,
		Password:         pharmacistReq.Password,
	}
}

type UpdatePharmacistConverter struct{}

func (c UpdatePharmacistConverter) ToEntity(pharmacistReq dto.UpdatePharmacistRequest) entity.Pharmacist {
	return entity.Pharmacist{
		PharmacyID:       pharmacistReq.PharmacyID,
		PhoneNumber:      pharmacistReq.PhoneNumber,
		YearOfExperience: pharmacistReq.YearOfExperience,
	}
}

type UpdatePharmacistPhotoConverter struct{}

func (c UpdatePharmacistPhotoConverter) ToEntity(pharmacistReq dto.UpdatePharmacistPhotoRequest) entity.Pharmacist {
	return entity.Pharmacist{
		ProfilePhoto: pharmacistReq.ProfilePhoto,
	}
}

type PaginationConverter struct{}

func (c PaginationConverter) ToDto(pagination entity.Pagination) dto.Pagination {
	return dto.Pagination{
		CurrentPage:     pagination.CurrentPage,
		TotalPage:       pagination.TotalPage,
		TotalPharmacist: pagination.TotalPharmacist,
	}
}

type GetPharmacistsConverter struct{}

func (c GetPharmacistsConverter) ToDto(pharmacist entity.Pharmacist) dto.GetPharmacistResponse {
	return dto.GetPharmacistResponse{
		ID:               pharmacist.ID,
		PharmacyID:       pharmacist.PharmacyID,
		PharmacyName:     pharmacist.PharmacyName,
		Name:             pharmacist.Name,
		SipaNumber:       pharmacist.SipaNumber,
		PhoneNumber:      pharmacist.PhoneNumber,
		YearOfExperience: pharmacist.YearOfExperience,
		Email:            pharmacist.Email,
	}
}

type PharmacistFileConverter struct{}

func (c *PharmacistFileConverter) ToEntity(fileDTO dto.FileRequest) entity.File {
	return entity.File{
		File: fileDTO.File,
	}
}

type GetRandomPassConverter struct{}

func (c GetRandomPassConverter) ToDto(randomPass entity.Pharmacist) dto.GetRandomPassResponse {
	return dto.GetRandomPassResponse{
		Password: randomPass.Password,
	}
}
