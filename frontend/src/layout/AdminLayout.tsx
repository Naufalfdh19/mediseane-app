import { Outlet } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

import AdminHeader from "../components/AdminHeader";
import NotFound from "../pages/NotFound";

import { GetAccessToken } from "../utils/localstorage";
import { TokenPayload } from "../types/token";

export default function AdminLayout() {
  const token = GetAccessToken();
  if (!token || jwtDecode<TokenPayload>(token).role !== "admin") {
    return <NotFound />;
  } 

  return (
    <div className="h-screen flex flex-col">
      <AdminHeader />
      <Outlet />
    </div>
  );
}
