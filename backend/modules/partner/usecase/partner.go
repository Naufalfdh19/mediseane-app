package usecase

import (
	"context"
	"encoding/json"
	"math"
	"montelukast/modules/partner/entity"
	queryparams "montelukast/modules/partner/query_params"
	"montelukast/modules/partner/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type PartnerUsecase interface {
	AddPartner(c context.Context, partner entity.Partner) error
	DeletePartner(c context.Context, partnerID int) error
	UpdatePartner(c context.Context, partner entity.Partner) error
	GetPartners(c context.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) (*entity.PartnerList, error)
	GetPartner(c context.Context, partnerID int) (*entity.Partner, error)
	UpdatePartnerFromConsumer(c context.Context, partner entity.Partner) error
}

type partnerUsecaseImpl struct {
	r        repository.PartnerRepo
	tr       transaction.TransactorRepoImpl
	rabbitMQ *amqp.Channel
}

func NewPartnerUsecase(rabbitMQ *amqp.Channel, r repository.PartnerRepo, tr transaction.TransactorRepoImpl) partnerUsecaseImpl {
	return partnerUsecaseImpl{
		r:        r,
		tr:       tr,
		rabbitMQ: rabbitMQ,
	}
}

func (u partnerUsecaseImpl) AddPartner(c context.Context, partner entity.Partner) error {
	err := apperror.CheckPartnerTime(partner.StartHour, partner.EndHour)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPartner, err, err)
	}

	isValid := apperror.IsYearFoundedValid(partner.YearFounded)
	if !isValid {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPartner, apperror.ErrInvalidYearFounded, apperror.ErrInvalidYearFounded)
	}
	err = apperror.CheckActiveDays(partner.ActiveDays)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPartner, err, err)
	}

	isPartnerExists, err := u.r.IsPartnerExistsByName(c, partner.Name)
	if err != nil {
		return err
	}
	if isPartnerExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPartner, apperror.ErrPartnerAlreadyExists, apperror.ErrPartnerAlreadyExists)
	}

	err = u.r.AddPartner(c, partner)
	if err != nil {
		return err
	}

	return nil
}

func calculateDelayUntilMidnight() int {
	now := time.Now()

	midnight := time.Date(
		now.Year(), now.Month(), now.Day()+1,
		0, 1, 0, 0, // 12:01 AM
		now.Location(),
	)

	delay := midnight.Sub(now).Milliseconds()
	return int(delay)
}

func (u partnerUsecaseImpl) UpdatePartner(c context.Context, partner entity.Partner) error {
	err := apperror.CheckPartnerTime(partner.StartHour, partner.EndHour)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPartner, err, err)
	}

	err = apperror.CheckActiveDays(partner.ActiveDays)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPartner, err, err)
	}

	isPartnerExists, err := u.r.IsPartnerExistsByID(c, partner.ID)
	if err != nil {
		return err
	}
	if !isPartnerExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrUpdatePartner, apperror.ErrPartnerNotExists, apperror.ErrPartnerNotExists)
	}
	delay := calculateDelayUntilMidnight()
	err = u.PublishDelayedMessage(c, partner, delay)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPartner, err, err)
	}
	return nil
}

func (u partnerUsecaseImpl) UpdatePartnerFromConsumer(c context.Context, partner entity.Partner) error {
	err := u.r.UpdatePartner(c, partner)
	if err != nil {
		return err
	}
	return nil
}

func (u partnerUsecaseImpl) PublishDelayedMessage(c context.Context, data entity.Partner, delay int) error {
	err := u.rabbitMQ.ExchangeDeclare(
		"update-partner-exchange", //name
		"x-delayed-message",       //type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		amqp.Table{
			"x-delayed-type": "fanout",
		},
	)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrUpdatePartner, apperror.ErrInternalServer, err)
	}
	body, err := json.Marshal(map[string]interface{}{
		"update_partner": data,
	})

	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrUpdatePartner, apperror.ErrInternalServer, err)
	}
	err = u.rabbitMQ.Publish(
		"update-partner-exchange",
		"x-delayed-message",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Headers: amqp.Table{
				"x-delay": delay,
			},
		},
	)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	return nil
}

func (u partnerUsecaseImpl) DeletePartner(c context.Context, partnerID int) error {
	isExists, err := u.r.IsPartnerExistsByID(c, partnerID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrDeletePartner, apperror.ErrPartnerNotExists, apperror.ErrPartnerNotExists)
	}

	isPharmacyExists, err := u.r.IsPharmacyExistsByPartnerID(c, partnerID)
	if err != nil {
		return err
	}
	if isPharmacyExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrDeletePartner, apperror.ErrPharmacyHasBeenCreated, apperror.ErrPharmacyHasBeenCreated)
	}

	err = u.r.DeletePartner(c, partnerID)
	if err != nil {
		return err
	}

	return nil
}

func (u partnerUsecaseImpl) GetPartners(c context.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) (*entity.PartnerList, error) {
	err := checkQueryParams(&queryParams, param)
	if err != nil {
		return nil, err
	}

	totalPartner, err := u.r.GetTotalPartners(c, queryParams)
	if err != nil {
		return nil, err
	}

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalPage := int(math.Ceil(float64(totalPartner) / float64(queryParams.Limit)))
	if totalPartner <= 0 {
		totalPage = 1
	}

	queryParams.Page = apperror.CheckCurrentPage(queryParams.Page, totalPage, queryParams.Limit)

	partners, err := u.r.GetPartners(c, queryParams, totalPartner)
	if err != nil {
		return nil, err
	}

	pagination := entity.Pagination{
		CurrentPage:     queryParams.Page,
		TotalPage:       totalPage,
		TotalPharmacist: totalPartner,
	}

	PartnerList := entity.PartnerList{
		Pagination: pagination,
		Partners:   partners,
	}

	return &PartnerList, nil
}

func checkQueryParams(queryParams *queryparams.QueryParams, param queryparams.QueryParamsExistence) error {
	var convertError error
	if param.IsLimitExists {
		limitInt, err := strconv.Atoi(param.Limit)
		if err != nil {
			convertError = err
		}
		queryParams.Limit = limitInt
	}
	if param.IsPageExists {
		pageInt, err := strconv.Atoi(param.Page)
		if err != nil {
			convertError = err
		}
		queryParams.Page = pageInt
	}
	if convertError != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrGetPartner, apperror.ErrConvertVariableType, nil)
	}

	if param.IsNameExists {
		queryParams.Name = "%" + param.Name + "%"
	}
	if param.IsYearFoundedExists {
		queryParams.YearFounded = "%" + param.YearFounded + "%"
	}
	if param.IsActiveDaysExists {
		queryParams.ActiveDays = "%" + param.ActiveDays + "%"
	}
	if param.IsStartHourExists {
		queryParams.StartHour = "%" + param.StartHour + "%"
	}
	if param.IsEndHourExists {
		queryParams.EndHour = "%" + param.EndHour + "%"
	}
	if param.IsSortByExsts {
		queryParams.SortBy = param.SortBy
	}
	if param.IsOrderExists {
		queryParams.Order = param.Order
	}

	if param.IsActiveExists {
		if param.IsActive == "true" {
			queryParams.IsActive = "true"
		} else {
			queryParams.IsActive = "false"
		}
	}

	return nil
}

func (u partnerUsecaseImpl) GetPartner(c context.Context, partnerID int) (*entity.Partner, error) {
	isPartnerExists, err := u.r.IsPartnerExistsByID(c, partnerID)
	if err != nil {
		return nil, err
	}
	if !isPartnerExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetPartner, apperror.ErrPartnerNotExists, apperror.ErrPartnerNotExists)
	}

	partner, err := u.r.GetPartner(c, partnerID)
	if err != nil {
		return nil, err
	}

	return partner, nil
}
