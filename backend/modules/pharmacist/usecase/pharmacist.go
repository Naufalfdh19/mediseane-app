package usecase

import (
	"context"
	"fmt"
	"math"
	"montelukast/modules/pharmacist/entity"
	queryparams "montelukast/modules/pharmacist/query_params"
	"montelukast/modules/pharmacist/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	imageuploader "montelukast/pkg/imageuploader"
	jwttoken "montelukast/pkg/jwt_token"
	"montelukast/pkg/transaction"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/resendlabs/resend-go"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

type PharmacistUsecase interface {
	Login(c context.Context, user entity.Pharmacist) (token string, err error)
	AddPharmacist(c context.Context, pharmacist entity.Pharmacist) error
	UpdatePharmacist(c context.Context, pharmacist entity.Pharmacist) error
	UpdatePharmacistPhoto(c context.Context, file entity.File, pharmacistID int) (string, error)
	DeletePharmacist(c context.Context, pharmacistID int) error
	GetPharmacistDetail(c context.Context, id int) (*entity.Pharmacist, error)
	GetPharmacists(c context.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) (*entity.PharmacistList, error)
	GetRandomPass(c context.Context) (*entity.Pharmacist, error)
}

type pharmacistUsecaseImpl struct {
	r  repository.PharmacistRepo
	tr transaction.TransactorRepoImpl
	rc *resend.Client
}

func NewPharmacistUsecase(r repository.PharmacistRepo, tr transaction.TransactorRepoImpl, rc *resend.Client) pharmacistUsecaseImpl {
	return pharmacistUsecaseImpl{
		r:  r,
		tr: tr,
		rc: rc,
	}
}

func (u pharmacistUsecaseImpl) Login(c context.Context, user entity.Pharmacist) (token string, err error) {
	jwt_token := jwttoken.NewJWT()
	exist, err := u.r.IsUserExistByEmail(c, user.Email)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidEmailOrPassword, err)
	}
	pharmacist, err := u.r.GetUserByEmail(c, user.Email)
	if err != nil {
		return "", err
	}
	if pharmacist.Role != appconstant.ROLE_PHARMACY {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidEmailOrPassword, err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(pharmacist.Password), []byte(user.Password)); err != nil {
		return "", apperror.NewErrStatusBadRequest(appconstant.FieldErrLogin, apperror.ErrInvalidEmailOrPassword, err)
	}
	pharmacistIdStr := strconv.Itoa(int(pharmacist.ID))
	token, err = jwt_token.GenerateJwtTokenForAuth(appconstant.JwtTokenAuthType, pharmacistIdStr, pharmacist.Role)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return token, nil
}

func (u pharmacistUsecaseImpl) AddPharmacist(c context.Context, pharmacist entity.Pharmacist) error {
	err := u.tr.WithinTransaction(c, func(txCtx context.Context) error {

		err := apperror.ValidateUsernameAndPassword(pharmacist.Name, pharmacist.Password)
		if err != nil {
			return apperror.NewErrStatusBadRequest(appconstant.FieldAddPharmacist, apperror.ErrInvalidEmailOrPassword, err)
		}

		isPhoneNumberValid := apperror.IsPhoneNumberValid(pharmacist.PhoneNumber)
		if !isPhoneNumberValid {
			return apperror.NewErrStatusBadRequest(appconstant.FieldAddPharmacist, apperror.ErrPhoneNumberNotValid, apperror.ErrPhoneNumberNotValid)
		}

		err = u.CheckPharmacist(txCtx, pharmacist)
		if err != nil {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(pharmacist.Password), bcrypt.DefaultCost)
		if err != nil {
			return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		nonHashPassword := pharmacist.Password
		pharmacist.Password = string(hash)

		err = u.r.AddPharmacist(txCtx, &pharmacist)
		if err != nil {
			return err
		}

		err = u.r.AddPharmacistDetail(txCtx, pharmacist)
		if err != nil {
			return err
		}

		pharmacist.Password = nonHashPassword
		// err = u.sendPharmacistAccountDetail(pharmacist)
		// if err != nil {
		// 	return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		// }

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u pharmacistUsecaseImpl) sendPharmacistAccountDetail(pharmacist entity.Pharmacist) error {
	params := &resend.SendEmailRequest{
		From:    "mediSEAne <no-reply@mediseane.store>",
		To:      []string{pharmacist.Email},
		Subject: "[mediSEAne] Pharmacist Account Details",
		Html: fmt.Sprintf(
			`<p style="font-size:3rem;font-weigth:bold;margin:0px">
				medi<span style="color:#008081">SEA</span>ne
			</p>
			<p style="font-weight:bold">
				All Your <span style="color:#008081">Healthcare</span> Needs at Your Fingertips
			</p>
			<hr>
			<p style="font-weight:bold;margin-bottom:40px">Hello, a pharmacist account have been created using this email!</p>
			<p>Bellow is the account details!</p>
			<ul style="margin-bottom:40px">
				<li>Name: <b>%s</b></li>
				<li>Password: <b>%s</b></li>
				<li>SIPA Number: <b>%s</b></li>
				<li>Phone Number: <b>%s</b></li>
				<li>Year Of Experience: <b>%d</b></li>
			</ul>
			<p>Click the button bellow to login!</p>
			<a href="%s%s" style="text-decoration:none !important;cursor:pointer !important;color:#FAFAFA !important;">
				<div style="display:inline;background-color:#008081;padding:10px 30px;border:none;border-radius:30px;color:#FAFAFA !important;font-weight:600">
					Login
				</div>
			</a>`,
			pharmacist.Name,
			pharmacist.Password,
			pharmacist.SipaNumber,
			pharmacist.PhoneNumber,
			pharmacist.YearOfExperience,
			os.Getenv("WEB_URL"),
			appconstant.PHARMACIST_LOGIN_PATH,
		),
	}

	_, err := u.rc.Emails.Send(params)
	if err != nil {
		return err
	}
	return nil
}

func (u pharmacistUsecaseImpl) CheckPharmacist(c context.Context, pharmacist entity.Pharmacist) error {
	isEmailExists, err := u.r.IsPharmacistExistsByEmail(c, pharmacist.Email)
	if err != nil {
		return err
	}
	if isEmailExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldAddPharmacist, apperror.ErrEmailAlreadyExists, apperror.ErrEmailAlreadyExists)
	}

	isSipaNumberExists, err := u.r.IsPharmacistExistsBySipaNumber(c, pharmacist.SipaNumber)
	if err != nil {
		return err
	}
	if isSipaNumberExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldAddPharmacist, apperror.ErrSipaNumberAlreadyExists, apperror.ErrSipaNumberAlreadyExists)
	}

	isPhoneNumberExists, err := u.r.IsPharmacistExistsByPhoneNumber(c, pharmacist.PhoneNumber)
	if err != nil {
		return err
	}
	if isPhoneNumberExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldAddPharmacist, apperror.ErrPhoneNumberAlreadyExists, apperror.ErrPhoneNumberAlreadyExists)
	}

	return nil
}

func (u pharmacistUsecaseImpl) GetPharmacistDetail(c context.Context, id int) (*entity.Pharmacist, error) {
	isPharmacistExists, err := u.r.IsPharmacistExistsByID(c, id)
	if err != nil {
		return nil, err
	}
	if !isPharmacistExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetPharmacist, apperror.ErrPharmacistNotExists, apperror.ErrPharmacistNotExists)
	}

	pharmacist, err := u.r.GetPharmacistByID(c, id)
	if err != nil {
		return nil, err
	}

	return pharmacist, nil
}

func (u pharmacistUsecaseImpl) UpdatePharmacist(c context.Context, pharmacist entity.Pharmacist) error {
	isPhoneNumberValid := apperror.IsPhoneNumberValid(pharmacist.PhoneNumber)
	if !isPhoneNumberValid {
		return apperror.NewErrStatusBadRequest(appconstant.FieldUpdatePharmacist, apperror.ErrPhoneNumberNotValid, apperror.ErrPhoneNumberNotValid)
	}

	isExists, err := u.r.IsPharmacistExistsByID(c, pharmacist.ID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldUpdatePharmacist, apperror.ErrPharmacistNotExists, apperror.ErrPharmacistNotExists)
	}

	if pharmacist.PharmacyID != nil {
		isExists, err = u.r.IsPharmacyExistsByID(c, *pharmacist.PharmacyID)
		if err != nil {
			return err
		}
		if !isExists {
			return apperror.NewErrStatusNotFound(appconstant.FieldUpdatePharmacist, apperror.ErrPharmacyNotExists, apperror.ErrPharmacyNotExists)
		}
	}

	pharmacyID, err := u.r.GetPharmacyIDByPharmacistID(c, pharmacist.ID)
	if err != nil {
		return err
	}
	if pharmacyID != nil {
		totalPharmacist, err := u.r.GetTotalPharmacistByPharmacyID(c, *pharmacyID)
		if err != nil {
			return err
		}
		if totalPharmacist <= 1 {
			return apperror.NewErrStatusBadRequest(appconstant.FieldUpdatePharmacist, apperror.ErrUpdatePharmacyForPharmacist, apperror.ErrUpdatePharmacyForPharmacist)
		}
	}

	isExists, err = u.r.IsPhoneNumberExists(c, pharmacist.ID, pharmacist.PhoneNumber)
	if err != nil {
		return err
	}
	if isExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldUpdatePharmacist, apperror.ErrPhoneNumberNotValid, apperror.ErrPhoneNumberNotValid)
	}

	err = u.r.UpdatePharmacist(c, pharmacist)
	if err != nil {
		return err
	}

	return nil

}

func (u pharmacistUsecaseImpl) UpdatePharmacistPhoto(c context.Context, file entity.File, pharmacistID int) (string, error) {
	isPharmacistExists, err := u.r.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return "", err
	}
	if !isPharmacistExists {
		return "", apperror.NewErrStatusNotFound(appconstant.FieldUpdatePharmacistPhoto, apperror.ErrPharmacistNotFound, apperror.ErrPharmacistNotFound)
	}

	validate := validator.New()
	err = validate.Struct(file)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	uploadUrl, err := imageuploader.ImageUploadHelper(file.File)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	err = u.r.UpdatePharmacistPhoto(c, uploadUrl, pharmacistID)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

func (u pharmacistUsecaseImpl) DeletePharmacist(c context.Context, pharmacistID int) error {
	err := u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		isPharmacistExists, err := u.r.IsPharmacistExistsByID(c, pharmacistID)
		if err != nil {
			return err
		}
		if !isPharmacistExists {
			return apperror.NewErrStatusNotFound(appconstant.FieldDeletePharmacist, apperror.ErrPharmacistNotFound, apperror.ErrPharmacistNotFound)
		}

		pharmacyID, err := u.r.GetPharmacyIDByPharmacistID(c, pharmacistID)
		if err != nil {
			return err
		}
		if pharmacyID != nil {
			return apperror.NewErrStatusBadRequest(appconstant.FieldDeletePharmacist, apperror.ErrPharmacistAssignToPharmacy, apperror.ErrPharmacistAssignToPharmacy)
		}

		err = u.r.DeletePharmacist(txCtx, pharmacistID)
		if err != nil {
			return err
		}
		err = u.r.DeletePharmacistDetail(txCtx, pharmacistID)
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

func (u pharmacistUsecaseImpl) GetPharmacists(c context.Context, queryParams queryparams.QueryParams, param queryparams.QueryParamsExistence) (*entity.PharmacistList, error) {

	err := checkQueryParams(&queryParams, param)
	if err != nil {
		return nil, err
	}

	queryParams = queryparams.DefaultQuery(queryParams)

	TotalPharmacist, err := u.r.GetTotalPharmacist(c, queryParams)
	if err != nil {
		return nil, err
	}

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalPage := int(math.Ceil(float64(TotalPharmacist) / float64(queryParams.Limit)))
	if TotalPharmacist <= 0 {
		totalPage = 1
	}

	queryParams.Page = apperror.CheckCurrentPage(queryParams.Page, totalPage, queryParams.Limit)

	pharmacists, err := u.r.GetPharmacists(c, queryParams, TotalPharmacist)
	if err != nil {
		return nil, err
	}

	pagination := entity.Pagination{
		CurrentPage:     queryParams.Page,
		TotalPage:       totalPage,
		TotalPharmacist: TotalPharmacist,
	}

	pharmacistsList := entity.PharmacistList{
		Pagination:  pagination,
		Pharmacists: pharmacists,
	}

	return &pharmacistsList, nil
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
	if param.IsYearOfExperienceExists {
		yearOfExperienceInt, err := strconv.Atoi(param.YearOfExperience)
		if err != nil {
			convertError = err
		}
		queryParams.YearOfExperience = yearOfExperienceInt
	}
	if convertError != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetUser, apperror.ErrConvertVariableType, nil)
		return err
	}

	if param.IsNameExists {
		queryParams.Name = "%" + param.Name + "%"
	}

	if param.IsEmailExists {
		queryParams.Email = "%" + param.Email + "%"
	}

	if param.IsSipaNumberExists {
		queryParams.SipaNumber = "%" + param.SipaNumber + "%"
	}

	if param.IsPhoneNumberExists {
		queryParams.PhoneNumber = "%" + param.PhoneNumber + "%"
	}

	if param.IsSortByExists {
		queryParams.SortBy = param.SortBy
	}
	if param.IsOrderExists {
		queryParams.Order = param.Order
	}

	return nil
}

func (u pharmacistUsecaseImpl) GetRandomPass(c context.Context) (*entity.Pharmacist, error) {
	var randomPass string
	var isPasswordValid bool
	var err error
	for !isPasswordValid {
		randomPass, err = password.Generate(10, 3, 1, false, false)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		isPasswordValid = apperror.IsPasswordValid(randomPass)
	}

	pharmacist := entity.Pharmacist{}
	pharmacist.Password = randomPass

	return &pharmacist, err
}
