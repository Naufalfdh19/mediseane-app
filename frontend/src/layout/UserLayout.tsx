import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

import Header from "../components/Header";
import Footer from "../components/Footer";

import useAxios from "../hooks/useAxios";
import { useSetUserState, useUserState } from "../store/user";
import { UserState } from "../store/userSlice";
import { GetAccessToken } from "../utils/localstorage";
import { Response } from "../types/response";
import { TokenPayload } from "../types/token";
import {
  API_METHOD_GET,
  API_USER_PROFILE,
  BASE_PATH,
  PATH_HOME,
} from "../const/const";
import PageLoader from "../components/PageLoader";

export default function UserLayout() {
  const setUserState = useSetUserState();
  const userState = useUserState();
  const token = GetAccessToken();
  const navigate = useNavigate();

  if (token && jwtDecode<TokenPayload>(token).role === "user") {
    useAxios<Response<UserState>>(
      API_USER_PROFILE,
      API_METHOD_GET,
      false,
      false,
      (data) => {
        setUserState(data.data);
      }
    );
  }

  useEffect(() => {
    if (location.pathname === BASE_PATH) {
      navigate(PATH_HOME);
    }
  }, [navigate]);

  if (userState.is_loading === true) {
    return (
      <div className="min-h-screen flex flex-col">
        <PageLoader />
      </div>
    );
  }

  return (
    <>
      <div className="min-h-screen flex flex-col">
        <Header />
        <Outlet />
      </div>
      <Footer />
    </>
  );
}
