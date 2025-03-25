import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router-dom";

import Header from "../components/Header";

import { PATH_AUTH, PATH_LOGIN } from "../const/const";

export default function AuthLayout() {
  const navigate = useNavigate();

  useEffect(() => {
    if (location.pathname === PATH_AUTH) {
      navigate(PATH_LOGIN);
    }
  }, [navigate]);

  return (
    <div className="h-screen flex flex-col">
      <Header type="auth"></Header>
      <Outlet />
    </div>
  );
}
