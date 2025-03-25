import React, { useState } from "react";

import Button from "../../../components/Button";
import PageLoader from "../../../components/PageLoader";
import SuccessBox from "../../../components/SuccessBox";
import ErrorBox from "../../../components/ErrorBox";
import Modal from "../../../components/Modal/Modal";
import Input from "../../../components/Input";
import ConfirmBox from "../../../components/ConfirmBox";

import { useSetUserUsername, useUserState } from "../../../store/user";
import useAxios from "../../../hooks/useAxios";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import { Response } from "../../../types/response";
import {
  API_EDIT_USERNAME,
  API_METHOD_PATCH,
  USERNAME_REGEX,
} from "../../../const/const";

interface UpdateUsernameProps {
  onClick: () => void;
}

export default function UpdateUsername({ onClick }: UpdateUsernameProps) {
  const [username, setUsername] = useState<string>();
  const [isUsernameValid, setIsUsernameValid] = useState<boolean>(true);
  const [showUpdateConfirm, setShowUpdateConfirm] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const userState = useUserState();
  const setUserUsername = useSetUserUsername();

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_EDIT_USERNAME,
    API_METHOD_PATCH
  );

  function handleUsernameChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(USERNAME_REGEX);
    setUsername(e.target.value);
    setIsUsernameValid(regex.test(e.target.value));
  }

  function handleSubmit() {
    fetchData({ username: username }).then((res) => {
      if (res && res.message && username) {
        setUserUsername(username);
      }
    });
    setShowResult(true);
    setShowUpdateConfirm(false);
  }

  return (
    <>
      {showUpdateConfirm ? (
        <ConfirmBox
          type="update"
          onYes={handleSubmit}
          onCancel={() => setShowUpdateConfirm(false)}
        />
      ) : isLoading ? (
        <Modal>
          <PageLoader />
        </Modal>
      ) : showResult ? (
        <>
          {error ? (
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
                onClick();
              }}
            >
              Username updated succesfully!
            </SuccessBox>
          )}
        </>
      ) : (
        <Modal onClose={onClick}>
          <div className="w-[70%] flex flex-col gap-[22px]">
            <p className="text-center text-2xl text-primary-black font-bold">
              Update Username
            </p>
            <p className="text-center text-primary-black">
              Enter your account's new username
            </p>
            <form className="flex flex-col gap-8">
              <div className="h-[70px] flex-col gap-2">
                <Input
                  type="text"
                  name="username"
                  placeholder="Username"
                  value={username}
                  valid={isUsernameValid}
                  onChange={handleUsernameChange}
                />
                {!isUsernameValid && (
                  <p className="text-sm text-primary-red pl-1 pt-1">
                    Username should be 5-12 alphanumeric characters
                  </p>
                )}
              </div>
              <div className="flex gap-4 justify-center">
                <Button
                  onClick={() => setShowUpdateConfirm(true)}
                  disabled={
                    !username || !isUsernameValid || username === userState.name
                  }
                >
                  Update
                </Button>
                <Button type="ghost" onClick={onClick} submit={false}>
                  Close
                </Button>
              </div>
            </form>
          </div>
        </Modal>
      )}
    </>
  );
}
