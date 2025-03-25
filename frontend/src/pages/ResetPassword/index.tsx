import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Eye, EyeOff } from "lucide-react";

import AuthFormContainer from "../../components/AuthFormContainer";
import ErrorCard from "../../components/ErrorCard";
import SuccessCard from "../../components/SuccessCard";
import PageLoader from "../../components/PageLoader";
import Input from "../../components/Input";
import Button from "../../components/Button";
import NotFound from "../NotFound";

import useAxios from "../../hooks/useAxios";
import {
  API_CHECK_RESET_TOKEN,
  API_METHOD_PATCH,
  API_RESET_PASSWORD,
  PASSWORD_REGEX,
  PATH_HOME,
  PATH_LOGIN,
  RESET_PASSWORD_TITLE,
  TOKEN_KEY,
} from "../../const/const";
import { Response } from "../../types/response";
import useTitle from "../../hooks/useTitle";

export default function ResetPassword() {
  const [password, setPassword] = useState<string>();
  const [isPasswordValid, setIsPasswordValid] = useState<boolean>(true);
  const [isPasswordShown, setIsPasswordShown] = useState<boolean>(false);
  const [isSuccess, setIsSuccess] = useState<boolean>(false);

  useTitle(RESET_PASSWORD_TITLE)

  const queryParams = new URLSearchParams(window.location.search);
  const token = queryParams.get(TOKEN_KEY);

  if (!token) {
    return <NotFound />;
  }

  const { error: tokenErr, isLoading: isCheckLoading } = useAxios<
    Response<undefined>
  >(`${API_CHECK_RESET_TOKEN}${token}`);

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_RESET_PASSWORD,
    API_METHOD_PATCH
  );

  const navigate = useNavigate();

  function handlePasswordChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(PASSWORD_REGEX);
    setPassword(e.target.value);
    setIsPasswordValid(regex.test(e.target.value));
  }

  function handlePasswordToggleClick() {
    setIsPasswordShown((prev) => !prev);
  }

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    fetchData({ token: token, new_password: password }).then((res) => {
      if (res) {
        setIsSuccess(true);
        setTimeout(() => {
          navigate(PATH_HOME);
          navigate(PATH_LOGIN);
        }, 5000);
      }
    });
  }

  if (tokenErr) {
    return <NotFound />;
  }

  if (isCheckLoading) {
    return <PageLoader />;
  }

  return (
    <AuthFormContainer>
      <h1 className="text-center text-2xl font-bold text-primary-black">
        Reset Your Password
      </h1>
      <p className="text-center text-primary-black">
        Enter your account's new password
      </p>
      <form onSubmit={handleSubmit}>
        {error && <ErrorCard errors={error} />}
        <div className="mt-3 flex flex-col gap-2">
          <Input
            type={isPasswordShown ? "text" : "password"}
            name="password"
            placeholder="Password"
            valid={isPasswordValid}
            onChange={handlePasswordChange}
            icon={isPasswordShown ? <Eye /> : <EyeOff />}
            onIconClick={handlePasswordToggleClick}
          />
          <p
            className={`text-sm ${
              isPasswordValid ? "text-primary-gray" : "text-primary-red"
            } pl-1`}
          >
            Minimum 8 characters with at least one uppercase, one lowercase, one
            special character, and a number
          </p>
        </div>
        <div className="mt-[20px]">
          <Button
            disabled={!password || !isPasswordValid || isLoading || isSuccess}
          >
            {!isLoading ? "Reset password" : "Loading..."}
          </Button>
        </div>
      </form>
      {isSuccess && (
        <SuccessCard message="Password updated successfully. Redirecting to Login Page..." />
      )}
    </AuthFormContainer>
  );
}
