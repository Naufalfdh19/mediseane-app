import { useState } from "react";

import AuthFormContainer from "../../components/AuthFormContainer";
import Input from "../../components/Input";
import Button from "../../components/Button";
import ErrorCard from "../../components/ErrorCard";
import ResetLinkSent from "./ResetLinkSent";

import {
  API_METHOD_POST,
  API_PATH_FORGET_PASSWORD,
  EMAIL_REGEX,
  FORGET_PASSWORD_TITLE,
} from "../../const/const";
import useAxios from "../../hooks/useAxios";
import { Response } from "../../types/response";
import useTitle from "../../hooks/useTitle";

export default function ForgetPassword() {
  const [email, setEmail] = useState<string>();
  const [isEmailValid, setIsEmailValid] = useState<boolean>(true);
  const [isSuccess, setIsSuccess] = useState<boolean>(false);

  useTitle(FORGET_PASSWORD_TITLE)

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_PATH_FORGET_PASSWORD,
    API_METHOD_POST
  );

  function handleEmailChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(EMAIL_REGEX);
    setEmail(e.target.value);
    setIsEmailValid(regex.test(e.target.value));
  }

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    fetchData({ email: email }).then((res) => {
      if (res && res.message) {
        setIsSuccess(true);
      }
    });
  }

  if (isSuccess) {
    return <ResetLinkSent />;
  }

  return (
    <AuthFormContainer>
      <h1 className="text-center text-2xl font-bold text-primary-black">
        Reset Your Password
      </h1>
      <p className="text-center text-primary-black">
        Enter your account's email address and we will send you a password reset
        link
      </p>
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
        <div className="mt-[20px]">
          <Button disabled={!email || !isEmailValid || isLoading}>
            {isLoading ? "Loading..." : "Send password reset link"}
          </Button>
        </div>
      </form>
    </AuthFormContainer>
  );
}
