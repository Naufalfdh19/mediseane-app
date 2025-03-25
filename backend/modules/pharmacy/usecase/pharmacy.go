package usecase

import (
	"context"
	"math"
	"montelukast/modules/pharmacy/entity"
	"montelukast/modules/pharmacy/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/imageuploader"

	"github.com/go-playground/validator"
)

type PharmacyUsecase interface {
	AddPharmacy(c context.Context, pharmacy entity.Pharmacy) (err error)
	UpdatePharmacy(c context.Context, pharmacy entity.Pharmacy) (err error)
	DeletePharmacy(c context.Context, id int) (err error)
	GetAllPharmacies(c context.Context, filter entity.PharmacyFilter) (pharmacies *entity.PaginatedPharmacies, err error)
	AddLogo(c context.Context, file entity.File) (string, error)
	GetPharmacyByID(c context.Context, id int) (pharmacy entity.Pharmacy, err error)
}

type pharmacyUsecaseImpl struct {
	pharmacyRepo repository.PharmacyRepository
}

func NewPharmacyUsecase(pharmacyRepo repository.PharmacyRepository) pharmacyUsecaseImpl {
	return pharmacyUsecaseImpl{
		pharmacyRepo: pharmacyRepo,
	}
}

func (u pharmacyUsecaseImpl) GetPharmacyByID(c context.Context, id int) (pharmacy entity.Pharmacy, err error) {
	exists, err := u.pharmacyRepo.IsPharmacyExists(c, id)
	if err != nil {
		return pharmacy, err
	}
	if !exists {
		return pharmacy, apperror.NewErrStatusBadRequest(appconstant.FieldErrPharmacy, apperror.ErrDataNotExists, err)
	}
	pharmacy, err = u.pharmacyRepo.GetPharmacyByID(c, id)
	if err != nil {
		return pharmacy, err
	}
	return pharmacy, nil
}

func (u pharmacyUsecaseImpl) GetAllPharmacies(c context.Context, filter entity.PharmacyFilter) (pharmacies *entity.PaginatedPharmacies, err error) {
	pharmacies, err = u.pharmacyRepo.GetAllPharmacies(c, filter)
	if err != nil {
		return nil, err
	}
	var pharmacyCount entity.PharmacyFilterCount
	pharmacyCount.City = filter.City
	pharmacyCount.Name = filter.Name
	totalItem, err := u.pharmacyRepo.GetTotalItem(c, pharmacyCount)
	if err != nil {
		return nil, err
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	totalPage := math.Ceil(float64(totalItem) / float64(filter.GetLimit()))
	pharmacies.Pagination.TotalItem = totalItem
	pharmacies.Pagination.CurrentPage = filter.Page
	pharmacies.Pagination.TotalPage = int(totalPage)
	return pharmacies, nil
}

func (u pharmacyUsecaseImpl) AddPharmacy(c context.Context, pharmacy entity.Pharmacy) (err error) {
	err = u.pharmacyRepo.AddPharmacy(c, pharmacy)
	if err != nil {
		return err
	}

	return nil
}

func (u pharmacyUsecaseImpl) isPharmacistAssigned(c context.Context, pharmacyID int) (err error) {
	exists, err := u.pharmacyRepo.IsPharmacistExists(c, pharmacyID)
	if err != nil {
		return err
	}
	if !exists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrPharmacy, apperror.ErrPharmacyCannotBeActivated, err)
	}
	return nil
}

func (u pharmacyUsecaseImpl) UpdatePharmacy(c context.Context, pharmacy entity.Pharmacy) (err error) {
	exists, err := u.pharmacyRepo.IsPharmacyExists(c, pharmacy.ID)
	if err != nil {
		return err
	}
	if !exists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrPharmacy, apperror.ErrPharmacyNotExists, err)
	}
	if pharmacy.IsActive{
		err = u.isPharmacistAssigned(c, pharmacy.ID)
	}
	if err != nil {
		return err
	}
	err = u.pharmacyRepo.UpdatePharmacy(c, pharmacy)
	if err != nil {
		return err
	}
	return nil
}

func (u pharmacyUsecaseImpl) DeletePharmacy(c context.Context, id int) (err error) {
	exists, err := u.pharmacyRepo.IsPharmacistExists(c, id)
	if err != nil {
		return err
	}
	if exists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrPharmacy, apperror.ErrPharmacistExist, err)
	}
	exists, err = u.pharmacyRepo.IsPharmacyExists(c, id)
	if err != nil {
		return err
	}
	if !exists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrPharmacies, apperror.ErrDataNotExists, err)
	}
	err = u.pharmacyRepo.DeletePharmacy(c, id)
	if err != nil {
		return err
	}
	return nil
}

func (u pharmacyUsecaseImpl) AddLogo(c context.Context, file entity.File) (string, error) {
	validate := validator.New()
	err := validate.Struct(file)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	uploadUrl, err := imageuploader.ImageUploadHelper(file.File)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	err = u.pharmacyRepo.AddLogo(c, uploadUrl, file.ID)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}
