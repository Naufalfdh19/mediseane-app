export interface Err {
  field: string;
  detail: string;
}

export interface Response<T> {
  message: string;
  data: T;
}

export interface Token {
  access_token: string;
}

export interface Province {
  id: string;
  name: string;
}

export interface City {
  id: string;
  name: string;
  latitude?: string;
  longitude?: string;
}

export interface District {
  id: string;
  name: string;
  latitude?: string;
  longitude?: string;
}

export interface SubDistrict {
  id: string;
  name: string;
  postal_codes?: string;
}

export type AddressTypes = Province | City | District | SubDistrict;

export interface CurrLocation {
  province_id: string;
  city_id: string;
  district_id?: string;
  subdistrict_id?: string;
}

export interface UserAddress {
  id: number;
  name: string;
  phone_number: string;
  address: string;
  province: string;
  city: string;
  district: string;
  sub_district: string;
  postal_code: string;
  is_active: boolean;
}

export interface UserAddressWithID {
  name: string;
  phone_number: string;
  address: string;
  province_id: number;
  city_id: number;
  district_id: number;
  sub_district_id: number;
  postal_code: string;
  longitude: string;
  latitude: string;
  is_active: boolean;
}

export interface Pagination {
  current_page: number;
  total_page: number;
  total_item: number;
}

export interface UserListItem {
  id: number;
  name: string;
  email: string;
  profile_photo: string;
  role: string;
}

export interface UserData {
  pagination: Pagination;
  user_list: UserListItem[] | null;
}

export interface PharmacyItemWithID {
  id: number;
  partner_id: number;
  partner_name: string;
  name: string;
  address: string;
  province_id: number;
  province: string;
  district_id: number;
  district: string;
  sub_district_id: number;
  sub_district: string;
  city_id: number;
  city: string;
  postal_code: number;
  is_active: boolean;
  longitude: string;
  latitude: string;
  logo: string;
  updated_at: string;
}

export interface PharmacyListItem {
  id: number;
  partner_id: number;
  partner_name: string;
  name: string;
  address: string;
  province_id: number;
  province: string;
  district_id: number;
  district: string;
  sub_district_id: number;
  sub_district: string;
  city_id: number;
  city: string;
  postal_code: number;
  longitude: string;
  latitude: string;
  logo: string;
  updated_at: string;
}

export interface PharmacyData {
  pagination: Pagination;
  pharmacies: PharmacyListItem[];
}

export interface Partner {
  id: number;
  name: string;
  year_founded: string;
  active_days: string;
  start_hour: string;
  end_hour: string;
  is_active: boolean;
}

export interface PartnersData {
  pagination: Pagination;
  partner_list: Partner[] | null;
}

export interface ProductPagination {
  current_page: number;
  total_page: number;
  total_product: number;
}

export interface ProductListItem {
  id: number;
  pharmacy_product_id: number;
  image: string;
  name: string;
  manufacture: string;
  pharmacy_name: string;
  price: string;
}

export interface ProductResponse {
  pagination: ProductPagination;
  products: ProductListItem[];
}

export interface UserProductDetails {
  id: number;
  pharmacy_product_id: number;
  product_categories: string[];
  name: string;
  generic_name: string;
  manufacture: string;
  description: string;
  image: string;
  unit_in_pack: number;
  stock: number;
  price: string;
  address: string;
  pharmacies_name: string;
}

export interface CartOverviewListItem {
  cart_item_id: number;
  pharmacy_product_id: number;
  quantity: number;
}

export interface ProductPagination {
  current_page: number;
  total_page: number;
  total_product: number;
}

export interface ProductListItem {
  id: number;
  pharmacy_product_id: number;
  image: string;
  name: string;
  manufacture: string;
  pharmacy_name: string;
  price: string;
}

export interface ProductResponse {
  pagination: ProductPagination;
  products: ProductListItem[];
}

export interface UserProductDetails {
  id: number;
  pharmacy_product_id: number;
  product_categories: string[];
  name: string;
  generic_name: string;
  manufacture: string;
  description: string;
  image: string;
  unit_in_pack: number;
  stock: number;
  price: string;
  address: string;
  pharmacies_name: string;
}

export interface CartOverviewListItem {
  cart_item_id: number;
  pharmacy_product_id: number;
  quantity: number;
}

export interface UserOrderProductDetails {
  order_product_id: number;
  pharmacy_product_id: number;
  name: string;
  manufacturer: string;
  image: string;
  quantity: number;
  subtotal: string;
}

export interface UserOrderDetails {
  details_id: number;
  pharmacy_id: number;
  pharmacy_name: string;
  status: string;
  order_products: UserOrderProductDetails[];
  logistic_price: string;
}

export interface UserOrdersListItem {
  id: number;
  total_price: string;
  order_date: string;
  order_details: UserOrderDetails[];
}

export interface OrdersPagination {
  current_page: number;
  total_page: number;
  total_order: number;
}

export interface OrdersProductDetail {
  product_id: number;
  name: string;
  quantity: number;
  image: string;
}

export interface OrdersListItem {
  order_id: number;
  status: string;
  created_at: string;
  product_list?: OrdersProductDetail[] | null;
}

export interface OrdersData {
  pagination: OrdersPagination;
  orders: OrdersListItem[] | null;
}

export interface CartItem {
  cart_item_id: number;
  pharmacy_product_id: number;
  name: string;
  manufacturer: string;
  image: string;
  quantity: number;
  subtotal: string;
}

export interface GroupedCartItem {
  pharmacy_id: number;
  pharmacy_name: string;
  items: CartItem[];
}

export interface Pharmacist {
  id: number;
  name: string;
  sipa_number: string;
  phone_number: string;
  year_of_experience: number;
  email: string;
  pharmacy_id: number;
  pharmacy_name: string;
  created_at: string;
}

export interface PharmacistsData {
  pagination: Pagination;
  pharmacist_list: Pharmacist[];
}

export interface GeneratedPassword {
  password: string;
}

export interface CartItem {
  cart_item_id: number;
  pharmacy_product_id: number;
  name: string;
  manufacturer: string;
  image: string;
  quantity: number;
  subtotal: string;
}

export interface GroupedCartItem {
  pharmacy_id: number;
  pharmacy_name: string;
  items: CartItem[];
}

export interface CheckoutDetails {
  id: string;
  data: GroupedCartItem[];
}

export interface Delivery {
  id: number;
  name: string;
  cost: number;
  etd: string;
}

export interface DeliveryData {
  pharmacy_id: number;
  delivery_id: number;
}

export interface CheckoutData {
  id_cart: string;
  delivery_data_list: DeliveryData[];
}

export interface CategoryListItem {
  id: number;
  name: string;
  updated_at: string;
}

export interface CategoryData {
  pagination: Pagination;
  list_item: CategoryListItem[];
}

export interface PharmacyProduct {
  id: number;
  name: string;
  generic_name: string;
  manufacturer: string;
  product_classification: string;
  product_form: string;
  stock: number;
  is_active: boolean;
  price: string;
  created_at: string;
}

export interface PharmacyProductsData {
  pagination: Pagination;
  pharmacy_products: PharmacyProduct[];
}

export interface AdminProductListItem {
  id: number;
  product_classification: string;
  product_form: string;
  name: string;
  generic_name: string;
  manufacture: string;
  is_active: string;
}

export interface AdminProductData {
  pagination: ProductPagination;
  products: AdminProductListItem[];
}
