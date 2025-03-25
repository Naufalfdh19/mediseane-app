package setup

import (
	"database/sql"
	"montelukast/middleware"
	"time"

	categoryHandler "montelukast/modules/category/handler"
	categoryRepo "montelukast/modules/category/repository"
	categoryUsecase "montelukast/modules/category/usecase"

	productHandler "montelukast/modules/product/handler"
	productRepo "montelukast/modules/product/repository"
	productUsecase "montelukast/modules/product/usecase"

	pharmacyProductHandler "montelukast/modules/pharmacyproduct/handler"
	pharmacyProductRepo "montelukast/modules/pharmacyproduct/repository"
	pharmacyProductUsecase "montelukast/modules/pharmacyproduct/usecase"

	adminHandler "montelukast/modules/admin/handler"
	adminRepo "montelukast/modules/admin/repository"
	adminUsecase "montelukast/modules/admin/usecase"

	partnerHandler "montelukast/modules/partner/handler"
	partnerRepo "montelukast/modules/partner/repository"
	partUsecase "montelukast/modules/partner/usecase"

	pharmacyHandler "montelukast/modules/pharmacy/handler"
	pharmacyRepo "montelukast/modules/pharmacy/repository"
	pharmacyUsecase "montelukast/modules/pharmacy/usecase"

	addressHandler "montelukast/modules/address/handler"
	addressRepo "montelukast/modules/address/repository"
	addressUsecase "montelukast/modules/address/usecase"

	cartHandler "montelukast/modules/cart/handler"
	cartRepo "montelukast/modules/cart/repository"
	cartUsecase "montelukast/modules/cart/usecase"

	pharmacistHandler "montelukast/modules/pharmacist/handler"
	pharmacistRepo "montelukast/modules/pharmacist/repository"
	pharmacistUsecase "montelukast/modules/pharmacist/usecase"

	orderHandler "montelukast/modules/order/handler"
	orderRepo "montelukast/modules/order/repository"
	orderUsecase "montelukast/modules/order/usecase"

	deliveryHandler "montelukast/modules/delivery/handler"
	deliveryRepo "montelukast/modules/delivery/repository"
	deliveryUsecase "montelukast/modules/delivery/usecase"

	userOrderHandler "montelukast/modules/userorder/handler"
	userOrderRepo "montelukast/modules/userorder/repository"
	userorderUsecase "montelukast/modules/userorder/usecase"

	checkoutHandler "montelukast/modules/checkout/handler"
	checkoutRepo "montelukast/modules/checkout/repository"
	checkoutUsecase "montelukast/modules/checkout/usecase"

	"montelukast/modules/user/handler"
	"montelukast/modules/user/repository"
	"montelukast/modules/user/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"

	"montelukast/pkg/transaction"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	cors "github.com/itsjamie/gin-cors"
	"github.com/resendlabs/resend-go"
)

type Handler struct {
	UserHandler            handler.UserHandler
	CategoryHandler        categoryHandler.CategoryHandler
	ProductHandler         productHandler.ProductHandler
	AdminHandler           adminHandler.AdminHandler
	PartnerHandler         partnerHandler.PartnerHandler
	PharmacyHandler        pharmacyHandler.PharmacyHandler
	AddressHandler         addressHandler.AddressHandler
	CartHandler            cartHandler.CartHandler
	DeliveryHandler        deliveryHandler.DeliveryHandler
	PharmacistHandler      pharmacistHandler.PharmacistHandler
	OrderHandler           orderHandler.OrderHandler
	UserOrderHandler       userOrderHandler.UserOrderHandler
	CheckoutHandler        checkoutHandler.CheckoutHandler
	PharmacyProductHandler pharmacyProductHandler.PharmacyProductHandler
}

func SetUp(db *sql.DB, redisDB *redis.Client, resendClient *resend.Client, rabbitMQ *amqp.Channel) *gin.Engine {
	transaction := transaction.NewTransactorRepo(db)
	userRepository := repository.NewUserRepo(db)
	userUsecase := usecase.NewUserUsecase(userRepository, transaction, resendClient)
	userHandler := handler.NewUserHandler(userUsecase)

	productRepository := productRepo.NewProductRepo(db)
	productUsecase := productUsecase.NewProductUsecase(productRepository, transaction)
	productHandler := productHandler.NewProductHandler(productUsecase)

	pharmacyRepository := pharmacyRepo.NewPharmacyRepository(db)
	pharmacyUsecase := pharmacyUsecase.NewPharmacyUsecase(pharmacyRepository)
	pharmacyHandler := pharmacyHandler.NewPharmacyHandler(pharmacyUsecase)

	pharmacistRepository := pharmacistRepo.NewPharmacistsRepo(db)
	pharmacistUsecase := pharmacistUsecase.NewPharmacistUsecase(pharmacistRepository, transaction, resendClient)
	pharmacistHandler := pharmacistHandler.NewPharmacistHandler(pharmacistUsecase)

	pharmacyProductRepository := pharmacyProductRepo.NewPharmacyProductRepo(db, redisDB)
	pharmacyProductUsecase := pharmacyProductUsecase.NewPharmacyProductUsecase(pharmacyProductRepository, transaction, pharmacyRepository, productRepository, pharmacistRepository)
	pharmacyProductHandler := pharmacyProductHandler.NewPharmacyProductHandler(pharmacyProductUsecase)

	cartRepository := cartRepo.NewCartRepo(db, redisDB)
	cartUsecase := cartUsecase.NewCartUsecase(cartRepository, pharmacyProductRepository)
	cartHandler := cartHandler.NewCartHandler(cartUsecase)

	adminRepository := adminRepo.NewAdminRepository(db)
	adminUsecase := adminUsecase.NewAdminUsecase(adminRepository)
	adminHandler := adminHandler.NewAdminHandler(adminUsecase)

	partnerRepository := partnerRepo.NewPartnerRepo(db)
	partnerUsecase := partUsecase.NewPartnerUsecase(rabbitMQ, partnerRepository, transaction)
	partnerHandler := partnerHandler.NewPartnerHandler(partnerUsecase)

	addressRepository := addressRepo.NewAddressRepo(db)
	addressUsecase := addressUsecase.NewAddressUsecase(addressRepository, transaction)
	addressHandler := addressHandler.NewAddressHandler(addressUsecase)

	orderRepository := orderRepo.NewOrderRepo(db)
	orderusecase := orderUsecase.NewOrderUsecase(rabbitMQ, orderRepository, transaction)
	orderHandler := orderHandler.NewOrderHandler(orderusecase)

	categoryRepository := categoryRepo.NewCategoryRepo(db)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepository)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryUsecase)

	deliveryRepostiory := deliveryRepo.NewDeliveryRepository(db, redisDB)
	checkoutRepo := checkoutRepo.NewCheckoutRepo(db, redisDB)

	checkoutUsecase := checkoutUsecase.NewCheckoutUsecase(&checkoutRepo, deliveryRepostiory, transaction)
	checkoutHandler := checkoutHandler.NewCheckoutHandler(checkoutUsecase)

	deliveryUsecase := deliveryUsecase.NewDeliveryUsecase(deliveryRepostiory, &checkoutRepo)
	deliveryHandler := deliveryHandler.NewDeliveryHandler(&deliveryUsecase)

	userOrderRepostiory := userOrderRepo.NewUserOrderRepo(db)
	userOrderUsecase := userorderUsecase.NewUserOrderUsecase(rabbitMQ, userOrderRepostiory, transaction, userRepository)
	userOrderHandler := userOrderHandler.NewUserOrderHandler(userOrderUsecase)

	consumer := userorderUsecase.NewRabbitMQConsumer(rabbitMQ, userOrderUsecase)
	go consumer.ConsumeDelayedMessage()

	partnerConsumer := partUsecase.NewRabbitMQConsumerPartner(rabbitMQ, partnerUsecase)
	go partnerConsumer.ConsumeDelayedMessage()

	updateStatusConsumer := orderUsecase.NewRabbitMQConsumerStatus(rabbitMQ, orderusecase)
	go updateStatusConsumer.ConsumeDelayedMessage()

	router := SetRouter(Handler{
		UserHandler:            userHandler,
		CategoryHandler:        categoryHandler,
		ProductHandler:         productHandler,
		AdminHandler:           adminHandler,
		PartnerHandler:         partnerHandler,
		PharmacyHandler:        pharmacyHandler,
		AddressHandler:         addressHandler,
		CartHandler:            cartHandler,
		DeliveryHandler:        deliveryHandler,
		PharmacistHandler:      pharmacistHandler,
		OrderHandler:           orderHandler,
		UserOrderHandler:       userOrderHandler,
		CheckoutHandler:        checkoutHandler,
		PharmacyProductHandler: pharmacyProductHandler,
	})

	return router
}

func SetRouter(h Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ObserveMiddleware)
	r.Use(middleware.ErrorMiddleware)
	r.Use(middleware.LoggerMiddleware)
	r.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, PATCH, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge:         50 * time.Second,
		Credentials:    false,
	}))

	baseEndpoint := r.Group("/api/v1")

	/* GENERAL */

	userAuth := baseEndpoint.Group("/auth")
	userAuth.POST("/register", h.UserHandler.RegisterHandler)
	userAuth.POST("/login", h.UserHandler.LoginHandler)
	userAuth.POST("/forget-password", h.UserHandler.ForgetPasswordHandler)
	userAuth.PATCH("/reset-password", h.UserHandler.ResetPasswordHandler)
	userAuth.GET("/reset-password/check", h.UserHandler.CheckResetPassTokenHandler)
	userAuth.PATCH("/verify-email", h.UserHandler.VerifyEmailHandler)

	userGeneral := baseEndpoint.Group("/")
	userGeneral.GET("/products", middleware.CheckProductAuthorization, h.ProductHandler.GetUserProductsHandler)
	userGeneral.GET("/products/homepage", middleware.CheckProductAuthorization, h.ProductHandler.GetUserProductsHomepageHandler)
	userGeneral.GET("/products/:id", h.ProductHandler.GetProductDetailHandler)

	adminAuth := baseEndpoint.Group("/admin/auth")
	adminAuth.POST("/login", h.AdminHandler.Login)

	pharmacistAuth := baseEndpoint.Group("/pharmacists/auth")
	pharmacistAuth.POST("/login", h.PharmacistHandler.Login)

	pharmacistGeneral := baseEndpoint.Group("/pharmacist")
	pharmacistGeneral.GET("products/master", h.ProductHandler.GetMasterProductsHandler)

	/* PROTECTED */

	protected := baseEndpoint.Group("/")
	protected.Use(middleware.CheckAuthorization)

	/* USER PROTECTED */

	userProtected := protected.Group("/")
	userProtected.Use(middleware.AuthUserMiddleware)

	userProtected.Use(middleware.CheckAuthorization)
	userProtected.Use(middleware.AuthUserMiddleware)
	userProtected.PUT("/carts", h.CartHandler.AddToCartHandler)
	userProtected.DELETE("/carts/:id", h.CartHandler.DeleteFromCartHandler)
	userProtected.GET("/carts", h.CartHandler.GetGroupedCartItemsHandler)
	userProtected.GET("/carts/overview", h.CartHandler.GetCartItemsHandler)
	userProtected.POST("/carts/checkout", h.CartHandler.GetSelectedCartItemsHandler)
	userProtected.PATCH("/order-details/:order_id/payment", h.UserOrderHandler.UpdatePaymentHandler)

	/* ADMIN PROTECTED */
	adminProtected := protected.Group("/admin")
	adminProtected.Use(middleware.AuthAdminMiddleware)

	adminProtected.GET("/users", h.AdminHandler.GetUsersHandler)

	adminProtected.GET("/pharmacies", h.PharmacyHandler.GetAllPharmaciesHandler)
	adminProtected.GET("pharmacies/:id", h.PharmacyHandler.GetPharmacyByID)
	adminProtected.POST("/pharmacies", h.PharmacyHandler.AddPharmacyHandler)
	adminProtected.PUT("/pharmacies", h.PharmacyHandler.UpdatePharmacyHandler)
	adminProtected.PATCH("/pharmacies/logo", h.PharmacyHandler.AddLogoHandler)
	adminProtected.DELETE("/pharmacies/:id", h.PharmacyHandler.DeletePharmacyHandler)

	adminProtected.GET("/partners", h.PartnerHandler.GetPartnersHandler)
	adminProtected.GET("/partners/:id", h.PartnerHandler.GetPartnerHandler)
	adminProtected.POST("/partners", h.PartnerHandler.AddPartnerHandler)
	adminProtected.PATCH("/partners/:id", h.PartnerHandler.UpdatePartnerHandler)
	adminProtected.DELETE("/partners/:id", h.PartnerHandler.DeletePartnerHandler)

	adminProtected.GET("/pharmacists", h.PharmacistHandler.GetPharmacistsHandler)
	adminProtected.POST("/pharmacists", h.PharmacistHandler.AddPharmacistHandler)
	adminProtected.PATCH("/pharmacists/:id", h.PharmacistHandler.UpdatePharmacistHandler)
	adminProtected.PATCH("/pharmacists/:id/profile-photo", h.PharmacistHandler.UpdatePharmacistPhotoHandler)
	adminProtected.DELETE("/pharmacists/:id", h.PharmacistHandler.DeletePharmacistHandler)
	adminProtected.GET("/pharmacists/:id", h.PharmacistHandler.GetPharmacistDetailHandler)
	adminProtected.GET("/random-password", h.PharmacistHandler.GetRandomPassHandler)

	adminProtected.POST("/categories", h.CategoryHandler.AddCategoryHandler)
	adminProtected.PUT("/categories", h.CategoryHandler.UpdateCategoryHandler)
	adminProtected.DELETE("/categories/:id", h.CategoryHandler.DeleteCategoryHandler)
	adminProtected.GET("/categories/:id", h.CategoryHandler.GetCategoryDetailHandler)
	userGeneral.GET("/categories", h.CategoryHandler.GetCategoriesHandler)

	adminProtected.POST("/products", h.ProductHandler.AddProductHandler)
	adminProtected.PATCH("/products/:id", h.ProductHandler.UpdateProductHandler)
	adminProtected.DELETE("/products/:id", h.ProductHandler.DeleteProductHandler)
	adminProtected.GET("/products", h.ProductHandler.GetProductsAdminHandler)

	/* USER PROTECTED */

	addressAuth := protected.Group("/addresses")
	addressAuth.GET("/check", h.AddressHandler.GetCurrentLocationHandler)
	addressAuth.GET("/provinces", h.AddressHandler.GetProvincesHandler)
	addressAuth.GET("/cities", h.AddressHandler.GetCitiesHandler)
	addressAuth.GET("/districts", h.AddressHandler.GetDistrictsHandler)
	addressAuth.GET("/sub-districts", h.AddressHandler.GetSubDistrictsHandler)
	addressAuth.GET("/user", middleware.AuthUserMiddleware, h.AddressHandler.GetUserAddressesHandler)
	addressAuth.POST("/user", middleware.AuthUserMiddleware, h.AddressHandler.AddUserAddressHandler)
	addressAuth.PUT("/user", middleware.AuthUserMiddleware, h.AddressHandler.UpdateUserAddressHandler)
	addressAuth.GET("/user/:id", middleware.AuthUserMiddleware, h.AddressHandler.GetUserAddressHandler)
	addressAuth.DELETE("/user/:id", middleware.AuthUserMiddleware, h.AddressHandler.DeleteUserAddressHandler)

	profileProtected := protected.Group("/profiles")
	profileProtected.Use(middleware.AuthUserMiddleware)
	profileProtected.GET("/user", h.UserHandler.GetProfileHandler)
	profileProtected.PATCH("/user/username", h.UserHandler.UpdateNameHandler)
	profileProtected.POST("/user/send-email", h.UserHandler.SendEmailFromProfileHandler)

	generalProtected := protected.Group("/")
	generalProtected.Use(middleware.AuthUserMiddleware)

	userProtected.GET("/orders", h.UserOrderHandler.GetDetailedOrdersHandler)
	userProtected.PATCH("/orders/:order-detail-id/completion", h.UserOrderHandler.ConfirmDeliveryHandler)
	userProtected.GET("/carts/checkout/delivery", h.DeliveryHandler.GetOngkirCost)
	userProtected.POST("/carts/checkout/order", h.CheckoutHandler.CheckoutCartHandler)
	userProtected.PATCH("/carts/checkout/cancel/:order-id", h.CheckoutHandler.CancelOrder)

	/* PHARMACIST PROTECTED */

	pharmacistProtected := protected.Group("/pharmacist")
	pharmacistProtected.Use(middleware.AuthPharmacistMiddleware)
	pharmacistProtected.GET("/orders", h.OrderHandler.GetOrdersHandler)
	pharmacistProtected.GET("/orders/:id", h.OrderHandler.GetOrderedProductsHandler)
	pharmacistProtected.PATCH("/orders/:id", h.OrderHandler.UpdateOrderStatusHandler)
	pharmacistProtected.DELETE("/orders/:id", h.OrderHandler.DeleteOrderHandler)

	pharmacistProtected.POST("/products", h.PharmacyProductHandler.AddPharmacyProductHandler)
	pharmacistProtected.PATCH("/products/:id", h.PharmacyProductHandler.UpdatePharmacyProductHandler)
	pharmacistProtected.DELETE("/products/:id", h.PharmacyProductHandler.DeletePharmacyProductHandler)
	pharmacistProtected.GET("/products", h.PharmacyProductHandler.GetPharmacyProductsHandler)
	pharmacistProtected.GET("/products/:id", h.PharmacyProductHandler.GetPharmacyProductDetailHandler)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.NoRoute(func(c *gin.Context) {
		err := apperror.NewErrStatusNotFound(appconstant.FieldErrNotFound, apperror.ErrNotFound, nil)
		c.Error(err)
	})

	return r
}
