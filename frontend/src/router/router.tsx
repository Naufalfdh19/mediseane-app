import { createBrowserRouter } from "react-router-dom";

import UserLayout from "../layout/UserLayout";
import AuthLayout from "../layout/AuthLayout";
import UserPrivateLayout from "../layout/UserPrivateLayout";
import AuthPrivateLayout from "../layout/AuthPrivateLayout";
import AdminLayout from "../layout/AdminLayout";
import ForgetPassword from "../pages/ForgetPassword";
import PharmacistLayout from "../layout/PharmacistLayout";
import PharmacistOrders from "../pages/Pharmacist/Order";
import Register from "../pages/Register";
import ResetPassword from "../pages/ResetPassword";
import NotFound from "../pages/NotFound";
import EmailVerification from "../pages/EmailVerification";
import Login from "../pages/Login";
import CreateAddress from "../pages/CreateAddress";
import Addresses from "../pages/Addresses";
import EditAddress from "../pages/EditAddress";
import AdminLogin from "../pages/AdminLogin";
import PharmacistLogin from "../pages/PharmacistLogin";
import GetUser from "../pages/Admin/GetUser";
import Profile from "../pages/Profile";
import AdminPartners from "../pages/Admin/Partners";
import AddPartner from "../pages/Admin/AddPartner";
import Cart from "../pages/Cart";
import Pharmacy from "../pages/Admin/Pharmacy";
import CreatePharmacy from "../pages/Admin/Pharmacy/CreatePharmacy";
import EditPharmacy from "../pages/Admin/Pharmacy/EditPharmacy";
import AdminPharmacists from "../pages/Admin/Pharmacists";
import AddPharmacist from "../pages/Admin/AddPharmacist";
import EditPharmacist from "../pages/Admin/EditPharmacist";
import Category from "../pages/Admin/Category";
import EditCategory from "../pages/Admin/Category/EditCategory";
import CreateCategory from "../pages/Admin/Category/CreateCategory";
import AdminLogout from "../pages/Admin/Logout";
import Checkout from "../pages/Checkout";
import Products from "../pages/Products";
import ProductDetails from "../pages/ProductDetails";
import PharmacistLogout from "../pages/Pharmacist/Logout";
import Orders from "../pages/Orders";
import Homepage from "../pages/Homepage";
import AddProduct from "../pages/Pharmacist/Product/AddProduct";

import {
  PATH_AUTH,
  PATH_FORGET_PASSWORD,
  PATH_LOGIN,
  NO_ROUTE,
  PATH_REGISTER,
  PATH_RESET_PASSWORD,
  PATH_VERIFY_EMAIL,
  PATH_USER,
  PATH_CREATE_ADDRESS,
  PATH_ADDRESS,
  PATH_PROFILE,
  PATH_EDIT_ADDRESS,
  PATH_ADMIN,
  PATH_ADMIN_USERS,
  PATH_ADMIN_LOGIN,
  PATH_ADMIN_PARTNERS,
  PATH_ADMIN_ADD_PARTNERS,
  PATH_ADMIN_PHARMACY,
  PATH_ADMIN_CREATE_PHARMACY,
  PATH_ADMIN_EDIT_PHARMACY,
  PATH_PRODUCTS,
  PATH_PRODUCT_DETAILS,
  PATH_ADMIN_PHARMACIST,
  PATH_ADMIN_ADD_PHARMACIST,
  PATH_ADMIN_EDIT_PHARMACIST,
  PATH_PHARMACIST_LOGIN,
  PATH_USER_CART,
  PATH_PHARMACIST,
  PATH_PHARMACIST_ORDERS,
  PATH_ADMIN_CATEGORY,
  PATH_ADMIN_CREATE_CATEGORY,
  PATH_ADMIN_EDIT_CATEGORY,
  PATH_ADMIN_DASHBOARD,
  PATH_ADMIN_LOGOUT,
  PATH_PHARMACIST_LOGOUT,
  PATH_PHARMACIST_DASHBOARD,
  PATH_CHECKOUT,
  PATH_USER_ORDERS,
  PATH_HOME,
  BASE_PATH,
  PATH_ADMIN_EDIT_PARTNERS,
  PATH_PHARMACIST_PRODUCTS,
  PATH_PHARMACIST_ADD_PRODUCT,
  PATH_PHARMACIST_EDIT_PRODUCT,
  PATH_ADMIN_PRODUCTS,
} from "../const/const";

import EditPartner from "../pages/Admin/EditPartner";
import Product from "../pages/Pharmacist/Product";
import EditProduct from "../pages/Pharmacist/Product/UpdateProduct";
import AdminProducts from "../pages/Admin/Products";
import PharmacistDashboard from "../pages/Pharmacist/Dashboard";
import AdminDashboard from "../pages/Admin/Dashboard/index";
import config from "../config";

function useRouter() {
  return createBrowserRouter([
    {
      path: NO_ROUTE,
      element: <NotFound />,
    },
    {
      path: BASE_PATH,
      element: <UserLayout />,
      children: [
        {
          path: PATH_USER,
          element: <UserPrivateLayout />,
          children: [
            {
              path: PATH_PROFILE,
              element: <Profile />,
            },
            {
              path: PATH_ADDRESS,
              element: <Addresses />,
            },
            {
              path: PATH_CREATE_ADDRESS,
              element: <CreateAddress />,
            },
            {
              path: PATH_EDIT_ADDRESS,
              element: <EditAddress />,
            },
            {
              path: PATH_USER_CART,
              element: <Cart />,
            },
            {
              path: PATH_CHECKOUT,
              element: <Checkout />,
            },
            {
              path: PATH_USER_ORDERS,
              element: <Orders />,
            },
          ],
        },
        {
          path: PATH_HOME,
          element: <Homepage />,
        },
        {
          path: PATH_PRODUCTS,
          element: <Products />,
        },
        {
          path: PATH_PRODUCT_DETAILS,
          element: <ProductDetails />,
        },
      ],
    },
    {
      path: PATH_AUTH,
      element: <AuthLayout />,
      children: [
        {
          path: PATH_VERIFY_EMAIL,
          element: <EmailVerification />,
        },
        {
          element: <AuthPrivateLayout />,
          children: [
            {
              path: PATH_LOGIN,
              element: <Login />,
            },
            {
              path: PATH_REGISTER,
              element: <Register />,
            },
            {
              path: PATH_FORGET_PASSWORD,
              element: <ForgetPassword />,
            },
            {
              path: PATH_RESET_PASSWORD,
              element: <ResetPassword />,
            },
            {
              path: PATH_ADMIN_LOGIN,
              element: <AdminLogin />,
            },
            {
              path: PATH_PHARMACIST_LOGIN,
              element: <PharmacistLogin />,
            },
          ],
        },
      ],
    },
    {
      path: PATH_ADMIN,
      element: <AdminLayout />,
      children: [
        {
          path: PATH_ADMIN_LOGOUT,
          element: <AdminLogout />,
        },
        {
          path: PATH_ADMIN_DASHBOARD,
          element: <AdminDashboard />,
        },
        {
          path: PATH_ADMIN_USERS,
          element: <GetUser />,
        },
        {
          path: PATH_ADMIN_PARTNERS,
          children: [
            {
              index: true,
              element: <AdminPartners />,
            },
            {
              path: PATH_ADMIN_ADD_PARTNERS,
              element: <AddPartner />,
            },
            {
              path: PATH_ADMIN_EDIT_PARTNERS,
              element: <EditPartner />,
            },
          ],
        },
        {
          path: PATH_ADMIN_PHARMACY,
          children: [
            {
              index: true,
              element: <Pharmacy />,
            },
            {
              path: PATH_ADMIN_CREATE_PHARMACY,
              element: <CreatePharmacy />,
            },
            {
              path: PATH_ADMIN_EDIT_PHARMACY,
              element: <EditPharmacy />,
            },
          ],
        },
        {
          path: PATH_ADMIN_PHARMACIST,
          children: [
            {
              index: true,
              element: <AdminPharmacists />,
            },
            {
              path: PATH_ADMIN_ADD_PHARMACIST,
              element: <AddPharmacist />,
            },
            {
              path: PATH_ADMIN_EDIT_PHARMACIST,
              element: <EditPharmacist />,
            },
          ],
        },
        {
          path: PATH_ADMIN_CATEGORY,
          children: [
            {
              index: true,
              element: <Category />,
            },
            {
              path: PATH_ADMIN_CREATE_CATEGORY,
              element: <CreateCategory />,
            },
            {
              path: PATH_ADMIN_EDIT_CATEGORY,
              element: <EditCategory />,
            },
          ],
        },
        {
          path: PATH_ADMIN_PRODUCTS,
          element: <AdminProducts />,
        },
      ],
    },
    {
      path: PATH_PHARMACIST,
      element: <PharmacistLayout />,
      children: [
        {
          path: PATH_PHARMACIST_LOGOUT,
          element: <PharmacistLogout />,
        },
        {
          path: PATH_PHARMACIST_DASHBOARD,
          element: <PharmacistDashboard />,
        },
        {
          path: PATH_PHARMACIST_ORDERS,
          element: <PharmacistOrders />,
        },
        {
          path: PATH_PHARMACIST_PRODUCTS,
          children: [
            {
              index: true,
              element: <Product />,
            },
            {
              path: PATH_PHARMACIST_ADD_PRODUCT,
              element: <AddProduct />,
            },
            {
              path: PATH_PHARMACIST_EDIT_PRODUCT,
              element: <EditProduct />,
            },
          ],
        },
      ],
    },
  ],
  {
    basename: config.ROUTE_BASE_PATH
  });
}

export default useRouter;
