import { Navigate, Outlet } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

import { GetAccessToken } from "../utils/localstorage";
import { TokenPayload } from "../types/token";
import { PATH_LOGIN } from "../const/const";


export default function UserPrivateLayout() {
  const token = GetAccessToken()

  if (!token || jwtDecode<TokenPayload>(token).role !== "user") {
    return <Navigate to={PATH_LOGIN} />;
  }

  return <Outlet />;
}
