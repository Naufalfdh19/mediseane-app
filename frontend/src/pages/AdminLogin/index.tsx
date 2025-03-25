import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Eye, EyeOff } from "lucide-react";

import AuthFormContainer from "../../components/AuthFormContainer";
import ErrorCard from "../../components/ErrorCard";
import Input from "../../components/Input";
import Button from "../../components/Button";

import useAxios from "../../hooks/useAxios";
import { SetAccessToken } from "../../utils/localstorage";
import { Response, Token } from "../../types/response";
import {
  ADMIN_LOGIN_TITLE,
  API_ADMIN_LOGIN,
  API_METHOD_POST,
  EMAIL_REGEX,
  PATH_ADMIN_DASHBOARD,
} from "../../const/const";
import useTitle from "../../hooks/useTitle";


export default function AdminLogin() {
  const [email, setEmail] = useState<string>();
  const [isEmailValid, setIsEmailValid] = useState<boolean>(true);
  const [password, setPassword] = useState<string>();
  const [isPasswordShown, setIsPasswordShown] = useState<boolean>(false);

  useTitle(ADMIN_LOGIN_TITLE)

  const navigate = useNavigate();

  const { error, isLoading, fetchData } = useAxios<Response<Token>>(
    API_ADMIN_LOGIN,
    API_METHOD_POST
  );

  function handleEmailChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(EMAIL_REGEX);
    setEmail(e.target.value);
    setIsEmailValid(regex.test(e.target.value));
  }

  function handlePasswordChange(e: React.ChangeEvent<HTMLInputElement>) {
    setPassword(e.target.value);
  }

  function handlePasswordToggleClick() {
    setIsPasswordShown((prev) => !prev);
  }

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    fetchData({ email: email, password: password }).then((res) => {
      if (res && res.data) {
        SetAccessToken(res.data.access_token);
        navigate(PATH_ADMIN_DASHBOARD);
      }
    });
  }

  return (
    <AuthFormContainer>
      <h1 className="text-center text-2xl font-bold text-primary-black">
        Admin Login
      </h1>
      <form onSubmit={handleSubmit}>
        {error && <ErrorCard errors={error} />}
        <div className="h-[70px] flex-col gap-2">
          <Input
            type="text"
            name="email"
            placeholder="Email"
            valid={isEmailValid}
            onChange={handleEmailChange}
          />
          {!isEmailValid && (
            <p className="text-sm text-primary-red pl-1 pt-1">
              Email is invalid
            </p>
          )}
        </div>
        <div className="mt-3 flex flex-col gap-2">
          <Input
            type={isPasswordShown ? "text" : "password"}
            name="password"
            placeholder="Password"
            onChange={handlePasswordChange}
            icon={isPasswordShown ? <Eye /> : <EyeOff />}
            onIconClick={handlePasswordToggleClick}
          />
        </div>
        <div className="mt-[20px]">
          <Button disabled={!email || !password || !isEmailValid || isLoading}>
            Login
          </Button>
        </div>
      </form>
    </AuthFormContainer>
  );
}
