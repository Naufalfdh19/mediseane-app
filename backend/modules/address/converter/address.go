package converter

import (
	"montelukast/modules/address/dto"
	"montelukast/modules/address/entity"
)

type GetProvincesConverter struct{}

func (c GetProvincesConverter) ToDto(location entity.Location) dto.ProvinceResponse {
	return dto.ProvinceResponse{
		ID:   location.ID,
		Name: location.Name,
	}
}

type GetCitiesAndDistrictsConverter struct{}

func (c GetCitiesAndDistrictsConverter) ToDto(location entity.Location) dto.CityAndDistrictResponse {
	return dto.CityAndDistrictResponse{
		ID:        location.ID,
		Name:      location.Name,
		Longitude: location.Longitude,
		Latitude:  location.Latitude,
	}
}

type GetSubDistrictsConverter struct{}

func (c GetSubDistrictsConverter) ToDto(location entity.Location) dto.SubDistrictResponse {
	return dto.SubDistrictResponse{
		ID:          location.ID,
		Name:        location.Name,
		PostalCodes: location.PostalCodes,
	}
}

type AddUserAddressConverter struct{}

func (c AddUserAddressConverter) ToEntity(userAddressReq dto.AddUserAddressRequest, userID int) entity.UserAddress {
	return entity.UserAddress{
		UserID:        userID,
		Name:          userAddressReq.Name,
		PhoneNumber:   userAddressReq.PhoneNumber,
		Address:       userAddressReq.Address,
		ProvinceID:    userAddressReq.ProvinceID,
		Province:      userAddressReq.Province,
		CityID:        userAddressReq.CityID,
		City:          userAddressReq.City,
		DistrictID:    userAddressReq.DistrictID,
		District:      userAddressReq.District,
		SubDistrictID: userAddressReq.SubDistrictID,
		SubDistrict:   userAddressReq.SubDistrict,
		PostalCode:    userAddressReq.PostalCode,
		Longitude:     userAddressReq.Longitude,
		Latitude:      userAddressReq.Latitude,
	}
}

type GetCurrentLocationConverter struct{}

func (c GetCurrentLocationConverter) ToDto(address entity.UserAddress) dto.GetCurrentLocationResponse {
	return dto.GetCurrentLocationResponse{
		ProvinceId: address.ProvinceID,
		CityId:     address.CityID,
	}
}

type GetUserAddressesConverter struct{}

func (c GetUserAddressesConverter) ToDto(address entity.UserAddress) dto.GetUserAddressesResponse {
	return dto.GetUserAddressesResponse{
		ID:          address.ID,
		Name:        address.Name,
		PhoneNumber: address.PhoneNumber,
		Address:     address.Address,
		Province:    address.Province,
		City:        address.City,
		District:    address.District,
		SubDistrict: address.SubDistrict,
		PostalCode:  address.PostalCode,
		IsActive:    address.IsActive,
	}
}

type UpdateUserAddressConverter struct{}

func (c UpdateUserAddressConverter) ToDto(addressDto dto.UpdateUserAddressRequest, userID int) entity.UserAddress {
	return entity.UserAddress{
		ID:            addressDto.ID,
		UserID:        userID,
		Name:          addressDto.Name,
		PhoneNumber:   addressDto.PhoneNumber,
		Address:       addressDto.Address,
		ProvinceID:    addressDto.ProvinceID,
		Province:      addressDto.Province,
		CityID:        addressDto.CityID,
		City:          addressDto.City,
		DistrictID:    addressDto.DistrictID,
		District:      addressDto.District,
		SubDistrictID: addressDto.SubDistrictID,
		SubDistrict:   addressDto.SubDistrict,
		PostalCode:    addressDto.PostalCode,
		Longitude:     addressDto.Longitude,
		Latitude:      addressDto.Latitude,
		IsActive:      *addressDto.IsActive,
	}
}

type GetUserAddressConverter struct{}

func (c GetUserAddressConverter) ToDto(address entity.UserAddress) dto.GetUserAddressResponse {
	return dto.GetUserAddressResponse{
		Name:          address.Name,
		PhoneNumber:   address.PhoneNumber,
		Address:       address.Address,
		ProvinceID:    address.ProvinceID,
		CityID:        address.CityID,
		DistrictID:    address.DistrictID,
		SubDistrictID: address.SubDistrictID,
		PostalCode:    address.PostalCode,
		Longitude:     address.Longitude,
		Latitude:      address.Latitude,
		IsActive:      address.IsActive,
	}
}
