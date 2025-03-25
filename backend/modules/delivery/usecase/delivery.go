package usecase

import (
	"context"
	"fmt"
	checkoutRepo "montelukast/modules/checkout/repository"
	"montelukast/modules/delivery/entity"
	"montelukast/modules/delivery/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

type DeliveryUsecase interface {
	GetAllOngkir(c context.Context, userID int, pharmacyID int) (ongkirList []entity.OngkirData, err error)
}

type deliveryUsecaseImpl struct {
	d repository.DeliveryRepository
	c checkoutRepo.CheckoutRepo
}

func NewDeliveryUsecase(d repository.DeliveryRepository, c checkoutRepo.CheckoutRepo) deliveryUsecaseImpl {
	return deliveryUsecaseImpl{
		d: d,
		c: c,
	}
}

func (u *deliveryUsecaseImpl) CalculateOngkirNextDay(c context.Context, req entity.OngkirRequest, userID int) (result entity.UserCostResponse, err error) {
	req.Weight = appconstant.MedicineWeight
	if req.LogisticPartnerID != appconstant.IDLogisticNextDay {
		return result, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer)
	}
	var calculateOngkir entity.CalculateOngkir
	exists, err := u.d.IsPostalCodeExist(c, req.DestinationPostalCode)
	if err != nil {
		return result, err
	}
	if !exists {
		calculateOngkir.DestinationID, err = u.d.GetOngkirLocationID(c, req.DestinationPostalCode)
		if err != nil {
			return result, err
		}
		err := u.d.AddLocationID(c, req.DestinationPostalCode, calculateOngkir.DestinationID)
		if err != nil {
			return result, err
		}
	} else {
		calculateOngkir.DestinationID, err = u.d.GetLocationID(c, fmt.Sprint(req.DestinationPostalCode))
		if err != nil {
			return result, err
		}
	}
	exists, err = u.d.IsPostalCodeExist(c, req.OriginPostalCode)
	if err != nil {
		return result, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	if !exists {
		calculateOngkir.OriginID, err = u.d.GetOngkirLocationID(c, req.OriginPostalCode)
		if err != nil {
			return result, err

		}
		err := u.d.AddLocationID(c, req.OriginPostalCode, calculateOngkir.OriginID)
		if err != nil {
			return result, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
	} else {
		calculateOngkir.OriginID, err = u.d.GetLocationID(c, fmt.Sprint(req.OriginPostalCode))
		if err != nil {
			return result, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
	}
	calculateOngkir.Weight = req.Weight
	calculateOngkir.Courier = appconstant.DefaultCourier
	calculateOngkir.BaseUrl = appconstant.URLOngkirCost
	calculateOngkir.SortingPrice = appconstant.LowestPrice
	result, err = u.d.GetOngkirCost(c, calculateOngkir)
	if err != nil {
		return result, err
	}
	return result, err
}

func (u *deliveryUsecaseImpl) GetAllOngkirRedis(c context.Context, addressID int, pharmacyID int) (ongkirList []entity.OngkirData, err error) {
	result, err := u.d.GetListOngkir(c, addressID, pharmacyID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *deliveryUsecaseImpl) StoreListOngkir(c context.Context, addressID int, pharmacyID int, ongkirList []entity.OngkirData) (err error) {
	err = u.d.StoreListOngkir(c, ongkirList, addressID, pharmacyID)
	if err != nil {
		return err
	}
	return nil
}

func (u *deliveryUsecaseImpl) GetAllOngkir(c context.Context, userID int, pharmacyID int) (ongkirList []entity.OngkirData, err error) {
	userPostal, addressID, err := u.d.GetUserPostalCode(c, userID)
	if err != nil {
		return nil, err
	}
	ongkirList, err = u.GetAllOngkirRedis(c, addressID, pharmacyID)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == nil {
		return ongkirList, nil
	}
	distance, err := u.d.CalculateDistance(c, pharmacyID, addressID)
	if err != nil {
		return nil, err
	}
	// if distance > 25000 {
	// 	return nil, apperror.NewErrStatusBadRequest(appconstant.FieldErrCheckout, apperror.ErrAddressTooFarr, apperror.ErrAddressTooFarr)
	// }
	sameDayPrice, err := u.d.GetSameDayPrice(c)
	if err != nil {
		return nil, err
	}
	instantPrice, err := u.d.GetInstantPrice(c)
	if err != nil {
		return nil, err
	}
	distanceInKM := distance / 1000
	calculateSameDay := sameDayPrice * distanceInKM
	calculateInstant := instantPrice * distanceInKM
	instantData := entity.OngkirData{
		Id:   appconstant.IDLogisticPartnerInstantDay,
		Name: appconstant.Instant,
		Cost: decimal.NewFromFloat(calculateInstant).Round(0),
		Etd:  appconstant.EtdInstant,
	}
	sameDayData := entity.OngkirData{
		Id:   appconstant.IDLogisticPartnerSameDay,
		Name: appconstant.SameDay,
		Cost: decimal.NewFromFloat(calculateSameDay).Round(0),
		Etd:  appconstant.EtdSameDay,
	}
	ongkirList = append(ongkirList, instantData, sameDayData)
	pharmacyPostal, err := u.d.GetPharmacyPostalCode(c, pharmacyID)
	if err != nil {
		return nil, err
	}
	data := entity.OngkirRequest{
		OriginPostalCode:      *pharmacyPostal,
		DestinationPostalCode: *userPostal,
		Weight:                appconstant.MedicineWeight,
		LogisticPartnerID:     appconstant.IDLogisticNextDay,
	}
	nextDayResult, err := u.CalculateOngkirNextDay(c, data, userID)
	if err == nil {
		nextDayData := entity.OngkirData{
			Id:   appconstant.IDLogisticNextDay,
			Name: nextDayResult.Name,
			Cost: nextDayResult.Cost.Round(1),
			Etd:  fmt.Sprint(nextDayResult.Etd),
		}
		ongkirList = append(ongkirList, nextDayData)
	}
	err = u.StoreListOngkir(c, addressID, pharmacyID, ongkirList)
	if err != nil {
		return nil, err
	}
	return ongkirList, nil
}
