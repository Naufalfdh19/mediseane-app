import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import { ChevronRight, LogOut, MailWarning } from "lucide-react";

import Button from "../../components/Button";
import TitleWithBackButton from "../../components/TitleWithBackButton";
import ConfirmBox from "../../components/ConfirmBox";
import Modal from "../../components/Modal/Modal";
import PageLoader from "../../components/PageLoader";
import ErrorBox from "../../components/ErrorBox";
import SuccessBox from "../../components/SuccessBox";
import UpdateUsername from "./UpdateUsername";

import useAxios from "../../hooks/useAxios";
import { useSetLoading, useUserState } from "../../store/user";
import { CapitalizeFirstLetter } from "../../utils/formatter";
import { RemoveAccessToken } from "../../utils/localstorage";
import {
  API_METHOD_POST,
  API_SEND_VERIFY_EMAIL,
  PATH_ADDRESS,
  PATH_HOME,
  PROFILE_TITLE,
} from "../../const/const";
import useTitle from "../../hooks/useTitle";

export default function Profile() {
  const [showLogoutConfirm, setShowLogoutConfirm] = useState<boolean>(false);
  const [showUpdateUsername, setShowUpdateUsername] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  useTitle(PROFILE_TITLE);

  const userState = useUserState();
  const setLoading = useSetLoading();

  const navigate = useNavigate();

  const { error, isLoading, fetchData } = useAxios(
    API_SEND_VERIFY_EMAIL,
    API_METHOD_POST
  );

  function handleLogout() {
    setLoading();
    navigate(PATH_HOME);
    setTimeout(() => {
      window.location.reload();
      RemoveAccessToken();
    }, 400);
  }

  function handleSendVerifyEmail() {
    fetchData({ email: userState.email });
    setShowResult(true);
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-7 md:my-16 w-[90%] max-w-[627px] flex flex-col gap-7 md:gap-10">
        <TitleWithBackButton>Profile</TitleWithBackButton>
        <div className="flex flex-col gap-4">
          <div className="w-full bg-white shadow-md p-10 border border-primary-gray/20 rounded-lg flex flex-col gap-6">
            <div className="flex justify-center">
              <img
                src={
                  userState.profile_photo ? userState.profile_photo : undefined
                }
                alt="user-profile"
                className="h-32 rounded-full"
              />
            </div>
            <div>
              <h1 className="text-primary-black text-center font-bold text-3xl">
                {userState.name}
              </h1>
              <p className="text-primary-gray text-center text-md md:text-xl">
                {userState.email}
              </p>
            </div>
            <div className="flex justify-center">
              <div className="w-full max-w-[300px]">
                <Button onClick={() => setShowUpdateUsername(true)}>
                  Edit
                </Button>
              </div>
            </div>
          </div>
          {!userState.is_verified && (
            <div className="w-full bg-white shadow-md p-6 border border-primary-gray/20 rounded-lg flex flex-col sm:flex-row justify-center sm:justify-between gap-2 items-center cursor-pointer">
              <div className="flex gap-2">
                <MailWarning className="text-primary-orange" />
                <p>Email isn't verified yet</p>
              </div>
              <div className="w-[130px]">
                <Button type="ghost" size="sm" onClick={handleSendVerifyEmail}>
                  Verify now
                </Button>
              </div>
            </div>
          )}
          <div
            className="w-full bg-white shadow-md p-6 border border-primary-gray/20 text-primary-blue rounded-lg flex justify-between cursor-pointer"
            onClick={() => navigate(PATH_ADDRESS)}
          >
            <p className="font-semibold">My addresses</p>
            <ChevronRight />
          </div>
          <div
            className="w-full bg-white shadow-md p-6 border border-primary-gray/20 text-primary-red rounded-lg flex justify-between cursor-pointer"
            onClick={() => setShowLogoutConfirm(true)}
          >
            <p className="font-semibold">Logout</p>
            <LogOut />
          </div>
        </div>
      </div>
      {showLogoutConfirm &&
        createPortal(
          <ConfirmBox
            type="logout"
            onYes={handleLogout}
            onCancel={() => setShowLogoutConfirm(false)}
          />,
          document.body
        )}
      {showUpdateUsername &&
        createPortal(
          <UpdateUsername onClick={() => setShowUpdateUsername(false)} />,
          document.body
        )}
      {showResult &&
        createPortal(
          <>
            {isLoading ? (
              <Modal>
                <PageLoader />
              </Modal>
            ) : error ? (
              <ErrorBox onClose={() => setShowResult(false)}>
                <p className="font-bold text-primary-black text-4xl">Oops...</p>
                <p className="text-primary-black text-xl">
                  {error && error[0].field === "server"
                    ? "Something is wrong, please try again"
                    : CapitalizeFirstLetter(error[0].detail)}
                </p>
              </ErrorBox>
            ) : (
              <SuccessBox
                onClose={() => {
                  setShowResult(false);
                }}
              >
                A verification link has been sent to your enail!
              </SuccessBox>
            )}
          </>,
          document.body
        )}
    </div>
  );
}
