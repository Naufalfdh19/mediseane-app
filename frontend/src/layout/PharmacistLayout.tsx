import { Outlet } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

import NotFound from "../pages/NotFound";
import PharmacistHeader from "../components/PharmacistHeader";

import { GetAccessToken } from "../utils/localstorage";
import { TokenPayload } from "../types/token";

export default function PharmacistLayout() {
  const token = GetAccessToken();
  if (!token || jwtDecode<TokenPayload>(token).role !== "pharmacist") {
    return <NotFound />;
  }

  return (
    <div className="h-screen flex flex-col">
      <PharmacistHeader />
      <Outlet />
    </div>
  );
}
