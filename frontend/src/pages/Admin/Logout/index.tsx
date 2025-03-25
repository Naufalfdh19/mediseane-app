import { useNavigate } from "react-router-dom";

import ConfirmBox from "../../../components/ConfirmBox";

import { RemoveAccessToken } from "../../../utils/localstorage";
import {
  ADMIN_LOGOUT_TITLE,
  PATH_ADMIN_LOGIN,
  PATH_BACK,
} from "../../../const/const";
import useTitle from "../../../hooks/useTitle";

export default function AdminLogout() {
  const navigate = useNavigate();

  useTitle(ADMIN_LOGOUT_TITLE);

  function handleLogout() {
    RemoveAccessToken();
    navigate(PATH_ADMIN_LOGIN);
  }

  return (
    <ConfirmBox
      type="logout"
      onYes={handleLogout}
      onCancel={() => navigate(PATH_BACK)}
    />
  );
}
