export const PATH_AUTH = "/auth";
export const PATH_LOGIN = "/auth/login";
export const PATH_REGISTER = "/auth/register";
export const PATH_FORGET_PASSWORD = "/auth/forget-password";
export const PATH_RESET_PASSWORD = "/auth/reset-password";
export const PATH_VERIFY_EMAIL = "/auth/verify-email";
export const PATH_ADMIN_LOGIN = "/auth/admin/login";
export const PATH_PHARMACIST_LOGIN = "/auth/pharmacist/login";

export const PATH_HOME = "/home";
export const PATH_USER = "/user";
export const PATH_PROFILE = "/user/profile";
export const PATH_ADDRESS = "/user/address";
export const PATH_CREATE_ADDRESS = "/user/create-address";
export const PATH_EDIT_ADDRESS_EMPTY = "/user/edit-address/";
export const PATH_EDIT_ADDRESS = PATH_EDIT_ADDRESS_EMPTY + ":id";
export const PATH_USER_CART = "/user/carts";
export const PATH_USER_ORDERS = "/user/orders";
export const PATH_CHECKOUT = "/user/checkout";
export const PATH_PRODUCTS = "/products";
export const PATH_PRODUCT_DETAILS_EMPTY = "/products/";
export const PATH_PRODUCT_DETAILS = PATH_PRODUCT_DETAILS_EMPTY + ":id";
export const BASE_PATH = "/";
export const PATH_BACK = -1;
export const NO_ROUTE = "*";

export const PATH_ADMIN = "/admin";
export const PATH_ADMIN_USERS = "/admin/user";
export const PATH_ADMIN_PHARMACY = "/admin/pharmacies";
export const PATH_ADMIN_CREATE_PHARMACY = "/admin/pharmacies/create";
export const PATH_ADMIN_EDIT_PHARMACY_EMPTY = "/admin/pharmacies/edit/";
export const PATH_ADMIN_EDIT_PHARMACY = PATH_ADMIN_EDIT_PHARMACY_EMPTY + ":id";
export const PATH_ADMIN_DASHBOARD = "/admin/dashboard";
export const PATH_ADMIN_PARTNERS = "/admin/partners";
export const PATH_ADMIN_ADD_PARTNERS = "/admin/partners/add";
export const PATH_ADMIN_EDIT_PARTNERS_EMPTY = "/admin/partners/edit/";
export const PATH_ADMIN_EDIT_PARTNERS = PATH_ADMIN_EDIT_PARTNERS_EMPTY + ":id";
export const PATH_ADMIN_PHARMACIST = "/admin/pharmacists";
export const PATH_ADMIN_ADD_PHARMACIST = "/admin/pharmacists/add";
export const PATH_ADMIN_EDIT_PHARMACIST_EMPTY = "/admin/pharmacists/edit/";
export const PATH_ADMIN_EDIT_PHARMACIST =
  PATH_ADMIN_EDIT_PHARMACIST_EMPTY + ":id";
export const PATH_ADMIN_CATEGORY = "/admin/category";
export const PATH_ADMIN_CREATE_CATEGORY = "/admin/category/create";
export const PATH_ADMIN_EDIT_CATEGORY_EMPTY = "/admin/category/edit/";
export const PATH_ADMIN_EDIT_CATEGORY = PATH_ADMIN_EDIT_CATEGORY_EMPTY + ":id";
export const PATH_ADMIN_PRODUCTS = "/admin/products";
export const PATH_ADMIN_LOGOUT = "/admin/logout";

export const PATH_PHARMACIST = "/pharmacist";
export const PATH_PHARMACIST_DASHBOARD = "/pharmacist/dashboard";
export const PATH_PHARMACIST_ORDERS = "/pharmacist/orders";
export const PATH_PHARMACIST_LOGOUT = "/pharmacist/logout";

export const PATH_PHARMACIST_PRODUCTS = "/pharmacist/products";
export const PATH_PHARMACIST_ADD_PRODUCT = "/pharmacist/products/add";
export const PATH_PHARMACIST_EDIT_PRODUCT_EMPTY = "/pharmacist/products/edit/";
export const PATH_PHARMACIST_EDIT_PRODUCT =
  PATH_PHARMACIST_EDIT_PRODUCT_EMPTY + ":id";

export const API_PATH_LOGIN = "/auth/login";
export const API_PATH_FORGET_PASSWORD = "/auth/forget-password";
export const API_PATH_REGISTER = "/auth/register";
export const API_CHECK_RESET_TOKEN = "/auth/reset-password/check?token=";
export const API_RESET_PASSWORD = "/auth/reset-password";
export const API_VERIFY_EMAIL = "/auth/verify-email";
export const API_ADMIN_LOGIN = "/admin/auth/login";
export const API_USER_PROFILE = "/profiles/user";
export const API_EDIT_USERNAME = "/profiles/user/username";
export const API_ADDRESS_USER = "/addresses/user";
export const API_ADDRESS_CHECK = "/addresses/check";
export const API_ADDRESS_PROVINCE = "/addresses/provinces";
export const API_ADDRESS_CITY = "/addresses/cities?province=";
export const API_ADDRESS_DISTRICT = "/addresses/districts?city=";
export const API_ADDRESS_SUBDISTRICT = "/addresses/sub-districts?district=";
export const API_USER_PRODUCTS = "/products";
export const API_GENERAL_PRODUCTS = "/general-products";
export const API_USER_HOMEPAGE = "/products/homepage";
export const API_GENERAL_HOMEPAGE = "/general-products/homepage";
export const API_USER_PRODUCT_DETAILS = "/products/";
export const API_USER_CARTS = "/carts";
export const API_USER_CARTS_OVERVIEW = "/carts/overview";
export const API_CART = "/carts";
export const API_USER_ORDERS = "/orders";
export const API_USER_ORDERS_SUFF_COMPLETION = "/completion";
export const API_USER_ORDER_DETAILS = "/order-details";
export const API_USER_ORDER_DETAILS_SUFF_PAYMENT = "/payment";
export const API_CHECKOUT = "/carts/checkout";
export const API_CHECKOUT_CONFIRM = "/carts/checkout/order";
export const API_CATEGORIES = "/categories";
export const API_SEND_VERIFY_EMAIL = "/profiles/user/send-email";
export const API_DELIVERY = "carts/checkout/delivery?pharmacy_id=";

export const API_ADMIN_USER = "/admin/users";
export const API_ADMIN_PARTNERS = "/admin/partners";
export const API_ADMIN_PHARMACY = "/admin/pharmacies";
export const API_ADMIN_PHARMACIST = "/admin/pharmacists";
export const API_ADMIN_GENERATED_PASSWORD = "/admin/random-password";
export const API_ADMIN_PRODUCT = "/admin/products";

export const API_MASTER_PRODUCTS = "/pharmacist/products/master";

export const API_PHARMACIST_LOGIN = "/pharmacists/auth/login";
export const API_PHARMACIST_ORDERS = "/pharmacist/orders";
export const API_ADMIN_CATEGORY = "/admin/categories";
export const API_PHARMACIST_PRODUCTS = "/pharmacist/products";

export const API_METHOD_POST = "POST";
export const API_METHOD_GET = "GET";
export const API_METHOD_PATCH = "PATCH";
export const API_METHOD_PUT = "PUT";
export const API_METHOD_DELETE = "DELETE";

export const TOKEN_KEY = "token";
export const BTN_HEIGHT_XS = "h-[28px]";
export const BTN_HEIGHT_SM = "h-[42px]";
export const BTN_HEIGHT_MD = "h-[60px]";
export const BTN_HEIGHT_LG = "h-[72px]";
export const BTN_HEIGHT_XL = "h-[90px]";

export const USERNAME_REGEX = /^[a-zA-Z0-9]{5,12}$/;
export const EMAIL_REGEX = /^[A-Za-z0-9._%+-]+@[A-Za-z0-9-]+[.][A-Za-z.]{2,}$/;
export const PASSWORD_REGEX =
  /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
export const PHONE_NUMBER_REGEX = /^(\+62|62|0)8[1-9][0-9]{6,9}$/;

export const ORDER_STATUS_WAITING = "waiting for payment";
export const ORDER_STATUS_PENDING = "pending";
export const ORDER_STATUS_PROCESSING = "processing";
export const ORDER_STATUS_SHIPPED = "shipped";
export const ORDER_STATUS_DELIVERED = "delivered";
export const ORDER_STATUS_COMPLETED = "completed";
export const ORDER_STATUS_CANCELED = "cancelled";

export const DAYS = [
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday",
];

export const TITLE_PREFIX = "mediSEAne - ";

export const HOMEPAGE_TITLE = "Home";
export const PRODUCTS_TITLE = "Products";

export const LOGIN_TITLE = "Login";
export const REGISTER_TITLE = "Register";
export const EMAIL_VERIFY_TITLE = "Email Verification";
export const FORGET_PASSWORD_TITLE = "Forget Password";
export const RESET_PASSWORD_TITLE = "Reset Password";

export const PROFILE_TITLE = "Profile";
export const ADDRESS_TITLE = "Address";
export const CREATE_ADDRESS_TITLE = "Create Address";
export const EDIT_ADDRESS_TITLE = "Edit Address";

export const CARTS_TITLE = "Carts";
export const CHECKOUT_TITLE = "Checkout";
export const ORDERS_TITLE = "Orders";

export const ADMIN_TITLE = "Admin | ";
export const ADMIN_LOGIN_TITLE = ADMIN_TITLE + "Login";
export const ADMIN_LOGOUT_TITLE = ADMIN_TITLE + "Logout";
export const ADMIN_DASHBOARD_TITLE = ADMIN_TITLE + "Dashboard";
export const ADMIN_USERS_TITLE = ADMIN_TITLE + "Users";
export const ADMIN_PRODUCTS_TITLE = ADMIN_TITLE + "Products";

export const ADMIN_PARTNERS_TITLE = ADMIN_TITLE + "Partners";
export const ADMIN_ADD_PARTNER_TITLE = ADMIN_TITLE + "Add Partner";
export const ADMIN_EDIT_PARTNER_TITLE = ADMIN_TITLE + "Edit Partner";

export const ADMIN_PHARMACISTS_TITLE = ADMIN_TITLE + "Pharmacists";
export const ADMIN_ADD_PHARMACIST_TITLE = ADMIN_TITLE + "Add Pharmacist";
export const ADMIN_EDIT_PHARMACIST_TITLE = ADMIN_TITLE + "Edit Pharmacist";

export const ADMIN_CATEGORY_TITLE = ADMIN_TITLE + "Category";
export const ADMIN_ADD_CATEGORY_TITLE = ADMIN_TITLE + "Add Category";
export const ADMIN_EDIT_CATEGORY_TITLE = ADMIN_TITLE + "Edit Category";

export const ADMIN_PHARMACIES_TITLE = ADMIN_TITLE + "Pharmacies";
export const ADMIN_ADD_PHARMACIES_TITLE = ADMIN_TITLE + "Add Pharmacy";
export const ADMIN_EDIT_PHARMACIES_TITLE = ADMIN_TITLE + "Edit Pharmacy";

export const PHARMACIST_TITLE = "Pharmacist | ";
export const PHARMACIST_LOGIN_TITLE = PHARMACIST_TITLE + "Login";
export const PHARMACIST_LOGOUT_TITLE = PHARMACIST_TITLE + "Logout";
export const PHARMACIST_DASHBOARD_TITLE = PHARMACIST_TITLE + "Dashboard";
export const PHARMACIST_ORDERS_TITLE = PHARMACIST_TITLE + "Orders";

export const PHARMACIST_PRODUCTS_TITLE = PHARMACIST_TITLE + "Products";
export const PHARMACIST_ADD_PRODUCT_TITLE = PHARMACIST_TITLE + "Add Product";
export const PHARMACIST_EDIT_PRODUCT_TITLE = PHARMACIST_TITLE + "Edit Product";
