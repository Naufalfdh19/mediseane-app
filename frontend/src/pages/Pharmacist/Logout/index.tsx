import { useNavigate } from "react-router-dom";

import ConfirmBox from "../../../components/ConfirmBox";

import { RemoveAccessToken } from "../../../utils/localstorage";
import { PATH_BACK, PATH_PHARMACIST_LOGIN, PHARMACIST_LOGOUT_TITLE } from "../../../const/const";
import useTitle from "../../../hooks/useTitle";

export default function PharmacistLogout() {
  const navigate = useNavigate();

  useTitle(PHARMACIST_LOGOUT_TITLE)

  function handleLogout() {
    RemoveAccessToken();
    navigate(PATH_PHARMACIST_LOGIN);
  }

  return (
    <ConfirmBox
      type="logout"
      onYes={handleLogout}
      onCancel={() => navigate(PATH_BACK)}
    />
  );
}
