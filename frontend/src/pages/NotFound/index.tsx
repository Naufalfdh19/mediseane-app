import { jwtDecode } from "jwt-decode";
import { useNavigate } from "react-router-dom";
import { SearchX } from "lucide-react";

import Button from "../../components/Button";

import { GetAccessToken } from "../../utils/localstorage";
import {
  PATH_ADMIN_DASHBOARD,
  PATH_HOME,
  PATH_PHARMACIST_DASHBOARD,
} from "../../const/const";
import { TokenPayload } from "../../types/token";

export default function NotFound() {
  const navigate = useNavigate();

  function renderButton() {
    const token = GetAccessToken();

    if (!token)
      return (
        <Button onClick={() => navigate(PATH_HOME)}>Return to home</Button>
      );

    const decode = jwtDecode<TokenPayload>(token);

    if (decode.role === "user")
      return (
        <Button onClick={() => navigate(PATH_HOME)}>Return to home</Button>
      );

    if (decode.role === "admin")
      return (
        <Button onClick={() => navigate(PATH_ADMIN_DASHBOARD)}>
          Return to dashboard
        </Button>
      );

    if (decode.role === "pharmacist")
      return (
        <Button onClick={() => navigate(PATH_PHARMACIST_DASHBOARD)}>
          Return to dashboard
        </Button>
      );
  }

  return (
    <div className="h-screen flex flex-col">
      <div className="grow flex flex-col justify-center items-center bg-primary-white">
        <div className="w-[90%] lg:w-[80%] mx-auto flex flex-col items-center gap-3">
          <SearchX
            className="h-32 w-32 md:h-64 md:w-64 text-primary-green"
            strokeWidth={0.5}
          />
          <h1 className="text-center text-2xl md:text-4xl font-bold text-primary-black">
            Page Not Found
          </h1>
          <p className="text-center md:text-xl text-primary-black">
            Sorry, the page you are looking for does not exist
          </p>
          <div className="w-[80%] md:w-[500px] md:mt-5">{renderButton()}</div>
        </div>
      </div>
    </div>
  );
}
