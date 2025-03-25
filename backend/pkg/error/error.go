package apperror

import (
	"errors"
	"net/http"
)

type Error struct {
	Field  string `json:"field" binding:"required"`
	Detail string `json:"detail" binding:"required"`
}

type ErrorStruct struct {
	Field         string
	Message       string
	Status        int
	SpecificError error
}

func (es ErrorStruct) Error() string {
	return es.SpecificError.Error()
}

func NewErrStatusBadRequest(field string, sentinel, err error) *ErrorStruct {
	return &ErrorStruct{
		Field:         field,
		Message:       sentinel.Error(),
		Status:        http.StatusBadRequest,
		SpecificError: err,
	}
}

func NewErrStatusUnauthorized(field string, sentinel, err error) *ErrorStruct {
	return &ErrorStruct{
		Field:         field,
		Message:       sentinel.Error(),
		Status:        http.StatusUnauthorized,
		SpecificError: err,
	}
}

func NewErrStatusNotFound(field string, sentinel, err error) *ErrorStruct {
	return &ErrorStruct{
		Field:         field,
		Message:       sentinel.Error(),
		Status:        http.StatusNotFound,
		SpecificError: err,
	}
}

func NewErrInternalServerError(field string, sentinel, err error) *ErrorStruct {
	return &ErrorStruct{
		Field:         field,
		Message:       sentinel.Error(),
		Status:        http.StatusInternalServerError,
		SpecificError: err,
	}
}

var (
	ErrDataNotExists               = errors.New("data not found")
	ErrIdEmpty                     = errors.New("id required ")
	ErrIdCity                      = errors.New("id should greater than 0")
	ErrPharmacistNotFound          = errors.New("pharmacist not found")
	ErrPhoneNumberNotValid         = errors.New("phone number is invalid")
	ErrSipaNumberAlreadyExists     = errors.New("sipa number already registered")
	ErrPhoneNumberAlreadyExists    = errors.New("phone number already registered")
	ErrPharmacistAlreadyExists     = errors.New("pharmacist already registered")
	ErrInvalidQuery                = errors.New("invalid query")
	ErrTokenIsInvalid              = errors.New("token is invalid")
	ErrorRegister                  = errors.New("error while register")
	ErrorUnique                    = errors.New("email has already been used")
	ErrInvalidAuth                 = errors.New("invalid authorization")
	ErrInsertData                  = errors.New("insert data unsuccessful")
	ErrEmptyId                     = errors.New("id must not be empty")
	ErrUpdateData                  = errors.New("error update data")
	ErrNoParamUpdate               = errors.New("no param for update")
	ErrTokenForgetPass             = errors.New("error while creating token")
	ErrResetPassword               = errors.New("error while resetting password")
	ErrPasswordInvalid             = errors.New("password must contain uppercase, lowercase, number, and special character")
	ErrDeleteData                  = errors.New("error delete data")
	ErrUpdatePharmacyForPharmacist = errors.New("pharmacist has less than 1 pharmacist assigned to them")
	ErrPharmacistAssignToPharmacy  = errors.New("pharmacist assign to pharmacy")
	ErrOrderDetailNotExists        = errors.New("order detail not exists")
	ErrProductNotExists            = errors.New("product not exists")
	ErrPharmacistNotExists         = errors.New("pharmacist not exists")
	ErrPharmacyNotExists           = errors.New("pharmacy not exists")
	ErrPharmacistUnauthorized      = errors.New("pharmacist is unauthorized")
	ErrInsert                      = errors.New("error insert into database")
	ErrUpdate                      = errors.New("error update database")
	ErrDelete                      = errors.New("error delete")
	ErrGetSubject                  = errors.New("error get subject from token")
	ErrEmailNotValid               = errors.New("email invalid")
	ErrPasswordNotValid            = errors.New("password invalid")
	ErrGetValueFromContext         = errors.New("value in context not exists")
	ErrConvertVariableType         = errors.New("convert variabel type failed")
	ErrParseJwtToken               = errors.New("parse jwt token failed")
	ErrQueryParams                 = errors.New("invalid query params")
	ErrPhoneNumber                 = errors.New("phone number is invalid")
	ErrAddressNotFound             = errors.New("address not found")
	ErrInvalidDeliveryData         = errors.New("invalid delivery data")
	ErrNoActiveAddress             = errors.New("at least one address need to be active")
	ErrStockUnavailable            = errors.New("product stock is unavailable")
	ErrCartNotAvailable            = errors.New("cart not exist")
	ErrPharmacyProductNotExists    = errors.New("pharmacy product does not exists")
	ErrCartItemNotExists           = errors.New("product not exists in cart")
	ErrInternalServer              = errors.New("internal server error")
	ErrUploadImage                 = errors.New("error uploading image")
	ErrUploadImageSize             = errors.New("image must smaller than 1MB")
	ErrIdType                      = errors.New("id must be integer")
	ErrNotFound                    = errors.New("sorry your destination is not found")
	ErrOrderCannotCanceled         = errors.New("order cannot be canceled")
	ErrCannotUpdateOrderStatus     = errors.New("cannot update order status")
	ErrUserNotExists               = errors.New("user not exists")
	ErrAddressNotExists            = errors.New("address not exists")
	ErrPartnerNotExists            = errors.New("partner not exists")
	ErrUnexpectedSigningMethod     = errors.New("unexpected signin method")
	ErrInvalidEmailOrPassword      = errors.New("email or password is invalid")
	ErrUsernameNotValid            = errors.New("username is invalid")
	ErrEmailAlreadyExists          = errors.New("email already registered")
	ErrPartnerAlreadyExists        = errors.New("partner already registered")
	ErrPharmacyHasBeenCreated      = errors.New("pharmacy has been created under this partner")
	ErrHashFailed                  = errors.New("hashing failed")
	ErrCredentialWrong             = errors.New("wrong credential")
	ErrTokenFailedToGenerated      = errors.New("token failed to generate")
	ErrInvalidJSON                 = errors.New("invalid JSON format")
	ErrInvalidDate                 = errors.New("invalid date(should use YYYY-MM-DD")
	ErrInvalidRangeDate            = errors.New("invalid range date")
	ErrTokenInvalid                = errors.New("token invalid or expired")
	ErrTokenNotExists              = errors.New("token not exists")
	ErrWrongLoginPage              = errors.New("this login page not match to your role")
	ErrSendEmail                   = errors.New("send email failed")
	ErrUserUnauthorized            = errors.New("user is unauthorized")
	ErrPageRange                   = errors.New("out of page range")
	ErrNoNearbyPharmacy            = errors.New("there is no nearby pharmacies in your location, please consider changing your location")
	ErrCategoryAlreadyExists       = errors.New("category already created")
	ErrProductClassNotExists       = errors.New("product classification not exists")
	ErrCategoryNotExists           = errors.New("product category not exists")
	ErrCategoryHasProductAlready   = errors.New("category has been assigned to product(s)")
	ErrProductFormMandatory        = errors.New("product form required")
	ErrUnitInPackMandatory         = errors.New("unit in pack required")
	ErrOrderNotExists              = errors.New("order not exists")
	ErrRajaOngkirNoResult          = errors.New("error raja ongkir no result")
	ErrPaymentAlreadyDone          = errors.New("payment already done")
	ErrOrderCannotBeCompleted      = errors.New("order cannot be completed yet")
	ErrAddressTooFarr              = errors.New("location too far")
	ErrInvalidOrderCancelation     = errors.New("invalid order")
	ErrInvalidDay                  = errors.New("invalid day")
	ErrDuplicateDay                = errors.New("duplicate day")
	ErrInvalidHour                 = errors.New("invalid hour")
	ErrInvalidYearFounded          = errors.New("invalid year founded")
	ErrWrongRole                   = errors.New("wrong role to login in this page")
	ErrFileEmpty                   = errors.New("file empty, please select a file to upload")
	ErrPharmacistNotHasPharmacy    = errors.New("pharmacist does not has pharmacy")
	ErrProductAlreadyExists        = errors.New("product already registered")
	ErrStockAlreadyUpdated         = errors.New("stock already updated")
	ErrPriceOrStockLessThanZero    = errors.New("product price or stock less than zero")
	ErrDuplicateCategory           = errors.New("duplicate category")
	ErrUserNotVerified             = errors.New("account is not verified yet")
	ErrPharmacyCannotBeActivated   = errors.New("pharmacy cannot be activated when there is no pharmacist yet")
	ErrMessageQueue                = errors.New("error publishing to message queue")
	ErrPharmacistExist             = errors.New("pharmacist still exists")
)
