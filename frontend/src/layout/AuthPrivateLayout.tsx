import { Navigate, Outlet } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

import { PATH_ADMIN_DASHBOARD, PATH_HOME } from "../const/const";
import { GetAccessToken } from "../utils/localstorage";
import { TokenPayload } from "../types/token";

export default function AuthPrivateLayout() {
  const token = GetAccessToken();
  if (token) {
    const decode = jwtDecode<TokenPayload>(token);
    switch (decode.role) {
      case "user":
        return <Navigate to={PATH_HOME} />;
      case "admin":
      case "pharmacist":
        return <Navigate to={PATH_ADMIN_DASHBOARD} />;
    }
  }

  return <Outlet />;
}
