import { useNavigate } from "react-router-dom";
import { FileText } from "lucide-react";
import { jwtDecode } from "jwt-decode";

import Input from "../Input";
import CartIcon from "../CartIcon";

import { useUserState } from "../../store/user";
import {
  API_METHOD_GET,
  API_USER_CARTS_OVERVIEW,
  PATH_HOME,
  PATH_LOGIN,
  PATH_PROFILE,
  PATH_REGISTER,
  PATH_USER_CART,
  PATH_USER_ORDERS,
} from "../../const/const";
import { GetAccessToken } from "../../utils/localstorage";
import { TokenPayload } from "../../types/token";
import useAxios from "../../hooks/useAxios";
import { CartOverviewListItem, Response } from "../../types/response";
import { useSetCartState } from "../../store/cart";

interface headerProps {
  type?: "default" | "auth";
  searchBar?: boolean;
}

function Header({ type = "default", searchBar = false }: headerProps) {
  const userState = useUserState();
  const navigate = useNavigate();

  const setCartState = useSetCartState();
  const token = GetAccessToken();

  if (token && jwtDecode<TokenPayload>(token).role === "user") {
    useAxios<Response<CartOverviewListItem[]>>(
      API_USER_CARTS_OVERVIEW,
      API_METHOD_GET,
      false,
      false,
      (data) => {
        setCartState(data.data);
      }
    );
  }

  return (
    <>
      <div className="min-h-14 w-full bg-primary-white flex items-center border-b border-solid border-primary-gray/30">
        <div className="w-[90%] lg:w-[80%] mx-auto flex justify-between">
          <h1
            className="text-base font-semibold text-primary-black cursor-pointer"
            onClick={() => {
              navigate(PATH_HOME);
            }}
          >
            medi<span className="text-primary-green">SEA</span>ne
          </h1>
          {type === "default" &&
            (userState.id ? (
              <div className="flex gap-4">
                <div
                  className="cursor-pointer"
                  onClick={() => {
                    navigate(PATH_USER_CART);
                  }}
                >
                  <CartIcon />
                </div>
                <FileText
                  className="cursor-pointer"
                  onClick={() => {
                    navigate(PATH_USER_ORDERS);
                  }}
                />
                <div
                  className="h-[24px] w-[24px] rounded-full text-center cursor-pointer"
                  onClick={() => {
                    navigate(PATH_PROFILE);
                  }}
                >
                  <img
                    src={
                      userState.profile_photo
                        ? userState.profile_photo
                        : undefined
                    }
                    alt="user-profile"
                    className="w-full rounded-full"
                  />
                </div>
              </div>
            ) : (
              <div className="text-primary-gray">
                <span
                  className="cursor-pointer hover:text-primary-green hover:font-bold"
                  onClick={() => navigate(PATH_LOGIN)}
                >
                  Login
                </span>{" "}
                /{" "}
                <span
                  className="cursor-pointer hover:text-primary-green hover:font-bold"
                  onClick={() => navigate(PATH_REGISTER)}
                >
                  Register
                </span>
              </div>
            ))}
        </div>
      </div>
      {searchBar && (
        <div>
          <Input type="text" name="search" placeholder="Search" />
        </div>
      )}
    </>
  );
}

export default Header;
