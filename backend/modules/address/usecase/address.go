package usecase

import (
	"context"
	"montelukast/modules/address/entity"
	"montelukast/modules/address/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
)

type AddressUsecase interface {
	GetProvinces(c context.Context) ([]entity.Location, error)
	GetCities(c context.Context, provinceID int) ([]entity.Location, error)
	GetDistricts(c context.Context, cityID int) ([]entity.Location, error)
	GetSubDistricts(c context.Context, districtID int) ([]entity.Location, error)
	AddUserAddress(c context.Context, userAddress entity.UserAddress) error
	GetCurrentLocation(c context.Context, longitude string, latittude string) (*entity.UserAddress, error)
	GetUserAddresses(c context.Context, userID int, filter entity.AddressFilter) ([]entity.UserAddress, error)
	UpdateUserAddress(c context.Context, address entity.UserAddress) error
	GetUserAddress(c context.Context, userID int, addressID int) (*entity.UserAddress, error)
	DeleteUserAddress(c context.Context, userID int, addressID int) error
}

type addressUsecaseImpl struct {
	r  repository.AddressRepo
	tr transaction.TransactorRepoImpl
}

func NewAddressUsecase(r repository.AddressRepo, tr transaction.TransactorRepoImpl) *addressUsecaseImpl {
	return &addressUsecaseImpl{
		r:  r,
		tr: tr,
	}
}

func (u *addressUsecaseImpl) GetProvinces(c context.Context) ([]entity.Location, error) {
	res, err := u.r.GetProvinces(c)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *addressUsecaseImpl) GetCities(c context.Context, provinceID int) ([]entity.Location, error) {
	res, err := u.r.GetCitiesByProvinceID(c, provinceID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *addressUsecaseImpl) GetDistricts(c context.Context, cityID int) ([]entity.Location, error) {
	res, err := u.r.GetDistrictsByCityID(c, cityID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *addressUsecaseImpl) GetSubDistricts(c context.Context, districtID int) ([]entity.Location, error) {
	res, err := u.r.GetSubDistrictsByDistrictID(c, districtID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *addressUsecaseImpl) AddUserAddress(c context.Context, userAddress entity.UserAddress) error {
	isPhoneNumberValid := apperror.IsPhoneNumberValid(userAddress.PhoneNumber)
	if !isPhoneNumberValid {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrPhoneNumber, apperror.ErrPhoneNumber, apperror.ErrPhoneNumber)
	}

	isExists, err := u.r.IsUserActiveAddressExists(c, userAddress.UserID)
	if err != nil {
		return err
	}
	userAddress.IsActive = !isExists

	err = u.r.AddUserAddress(c, userAddress)
	if err != nil {
		return err
	}

	return nil
}

func (u *addressUsecaseImpl) GetCurrentLocation(c context.Context, longitude string, latittude string) (*entity.UserAddress, error) {
	address, err := u.r.GetCurrentLocationByLongAndLat(c, longitude, latittude)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (u *addressUsecaseImpl) GetUserAddresses(c context.Context, userID int, filter entity.AddressFilter) ([]entity.UserAddress, error) {
	addresses, err := u.r.GetAddressesByUserID(c, userID, filter)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (u *addressUsecaseImpl) UpdateUserAddress(c context.Context, address entity.UserAddress) error {
	isExists, err := u.r.IsUserAddressExists(c, address.ID, address.UserID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrNotFound, apperror.ErrAddressNotFound, apperror.ErrAddressNotFound)
	}

	isPhoneNumberValid := apperror.IsPhoneNumberValid(address.PhoneNumber)
	if !isPhoneNumberValid {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrPhoneNumber, apperror.ErrPhoneNumber, apperror.ErrPhoneNumber)
	}

	activeID, err := u.r.GetActiveAddressByUserID(c, address.UserID)
	if err != nil {
		return err
	}
	if activeID == address.ID && !address.IsActive {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddress, apperror.ErrNoActiveAddress, apperror.ErrNoActiveAddress)
	}

	err = u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		if activeID != address.ID && address.IsActive {
			err = u.r.DeactivateAddressByID(txCtx, activeID)
			if err != nil {
				return err
			}
		}
		err = u.r.UpdateUserAddress(txCtx, address)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *addressUsecaseImpl) GetUserAddress(c context.Context, userID int, addressID int) (*entity.UserAddress, error) {
	isExists, err := u.r.IsUserAddressExists(c, addressID, userID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrNotFound, apperror.ErrAddressNotFound, apperror.ErrAddressNotFound)
	}

	address, err := u.r.GetAddressByID(c, addressID)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (u *addressUsecaseImpl) DeleteUserAddress(c context.Context, userID int, addressID int) error {
	isExists, err := u.r.IsUserAddressExists(c, addressID, userID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrNotFound, apperror.ErrAddressNotFound, apperror.ErrAddressNotFound)
	}

	activeID, err := u.r.GetActiveAddressByUserID(c, userID)
	if err != nil {
		return err
	}
	if activeID == addressID {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddress, apperror.ErrNoActiveAddress, apperror.ErrNoActiveAddress)
	}

	err = u.r.DeleteUserAddressByID(c, addressID)
	if err != nil {
		return err
	}

	return nil
}
