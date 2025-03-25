package appconstant

import "time"

const (
	CONFIG_SMTP_PORT      = 587
	CONFIG_SMTP_HOST      = "smtp.gmail.com"
	ROLE_USER             = "user"
	ROLE_ADMIN            = "admin"
	ROLE_PHARMACY         = "pharmacist"
	VERIFY_EMAIL_PATH     = "/auth/verify-email?token="
	RESET_PASSWORD_PATH   = "/auth/reset-password?token="
	PHARMACIST_LOGIN_PATH = "/auth/pharmacist/login"
	IMAGESIZEMAX          = 1 << 20
)

const (
	FieldErrOngkir                    = "get ongkir"
	FieldErrLogin                     = "login"
	FieldErrRegister                  = "register"
	FieldErrForgetPassword            = "forget password"
	FieldErrIsResetPassTokenExists    = "reset passord token"
	FieldErrIsVerifyEmailsTokenExists = "verify email token"
	FieldErrResetPassword             = "reset password"
	FieldErrVerifyEmail               = "verify email"
	FieldErrCheckAuthorization        = "authorization"
	FieldErrCheckUser                 = "check user"
	FieldErrCheckResetPassToken       = "check reset password token"
	FieldErrGetRole                   = "get role"
	FieldErrGetUser                   = "get user"
	FieldErrUpdateUser                = "update user"
	FieldErrGetListOfProducts         = "get list of products"
	FieldErrGetCities                 = "get cities"
	FieldErrGetDistricts              = "get districts"
	FieldErrGetSubDistricts           = "get sub-districts"
	FieldErrGetCurrentLocation        = "get current location"
	FieldErrGetAddresses              = "get addresses"
	FieldErrAddAddress                = "add address"
	FieldErrCheckout                  = "checkout"
	FieldErrPhoneNumber               = "phone number"
	FieldErrAddress                   = "address"
	FieldErrNotFound                  = "not found"
	FieldErrAddCategory               = "add product category"
	FieldErrUpdateCategory            = "update product category"
	FieldErrDeleteCategory            = "delete product category"
	FieldErrGetCategory               = "get category"
	FieldErrGetCategories             = "get list of categories"
	FieldErrJSON                      = "json"
	FieldErrGetPartners               = "get partner"
	FieldErrDeletePartner             = "delete partner"
	FieldErrServer                    = "server"
	FieldErrAddPartner                = "add partner"
	FieldErrUpdatePartner             = "update partner"
	FieldAddPharmacist                = "add pharmacist"
	FieldUpdatePharmacist             = "update pharmacist"
	FieldUpdatePharmacistPhoto        = "update pharmacist photo"
	FieldDeletePharmacist             = "delete pharmacist"
	FieldErrGetCart                   = "get cart"
	FieldErrGetPharmacist             = "get pharmacist"
	FieldAdminGetProducts             = "admin get products"
	FieldErrDeleteProduct             = "delete product"
	FieldErrGetPartner                = "get partner"
	FieldErrAddToCart                 = "add product to cart"
	FieldErrDeleteFromCart            = "delete product from cart"
	FieldErrGetPharmacies             = "retrieve pharmacies"
	FieldErrPharmacies                = "delete pharmacy"
	FieldErrPharmacy                  = "pharmacy"
	FieldErrUploadImage               = "error upload image"
	FieldErrImageType                 = "must upload image with extension png,jpg,or jpeg"
	FieldErrImageSize                 = "image size exceeded"
	FieldErrGetUserProduct            = "get user product"
	FieldErrUploadPayment             = "upload payment"
	FieldErrGetUserOrders             = "get user orders"
	FieldErrGetOrderedProducts        = "get ordered products"
	FieldErrDeleteOrder               = "delete order"
	FieldErrUpdateOrderStatus         = "update order status"
	FieldErrAddProduct                = "add product"
	FieldErrUpdateProduct             = "update product"
	FieldErrGetProducts               = "get products"
	FieldErrConfirmDelivery           = "confirm delivery"
	FieldErrGetProductDetail          = "get product details"
	FieldErrGetProduct                = "get product"
	FieldErrCancel                    = "cancel order"
	FieldErrEmail                     = "email"
	FieldErrAddPharmacyProduct        = "add pharmacy product"
	FieldErrUpdatePharmacyProduct     = "update pharmacy product"
	FieldErrGetPharmacyProducts       = "get pharmacy products"
	FieldUpdateProductPhoto           = "update product photo"
	FieldErrChangeStatus              = "error change status"
)

const (
	JwtTokenAuthType        = "auth"
	JwtTokenResetPassType   = "reset_password"
	JwtTokenVerifyEmailType = "verify email"
)

const (
	DefaultProfileIMG = "https://static.thenounproject.com/png/363639-200.png"
	DefaultLocation   = "POINT (106.74130492346 -6.191140471555)"
)

const (
	IDLogisticPartnerSameDay     = 2
	IDLogisticPartnerInstantDay  = 1
	IDLogisticNextDay            = 3
	SameDay                      = "Same Day"
	Instant                      = "Instant"
	EtdSameDay                   = "1 day"
	EtdInstant                   = "1 day"
	SameDayPrice                 = 1000000
	InstantPrice                 = 2500
	LimitInitialvaluePharmacist  = 10
	PageInitialvaluePharmacist   = 1
	SortbyInitialvaluePharmacist = "created_at"
	OrderInitialvaluePharmacist  = "desc"
	OngkirTimeExpiration         = 5 * time.Minute
	CartRedisExpiration          = 5 * time.Minute
	OngkirRedisKey               = "shipping:%d:%d"
	MedicineWeight               = 1000
	DefaultStatusOrder           = "Waiting for Payment"
	PaymentCompleteTime          = 60000
	OrderCompleteTime            = 10000
)

const (
	StatusCancelled  = "Cancelled"
	StatusDelivered  = "Delivered"
	StatusShipped    = "Shipped"
	StatusProcessing = "Processing"
	StatusPending    = "Pending"
)

const (
	URLOngkirLocationID = "https://rajaongkir.komerce.id/api/v1/destination/domestic-destination"
	URLOngkirCost       = "https://rajaongkir.komerce.id/api/v1/calculate/domestic-cost"
	DefaultCourier      = "jne"
	LowestPrice         = "lowest"
)
