import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Eye, EyeOff } from "lucide-react";

import AuthFormContainer from "../../components/AuthFormContainer";
import ErrorCard from "../../components/ErrorCard";
import Input from "../../components/Input";
import Button from "../../components/Button";

import useAxios from "../../hooks/useAxios";
import { Response } from "../../types/response";
import {
  API_METHOD_POST,
  API_PATH_REGISTER,
  EMAIL_REGEX,
  USERNAME_REGEX,
  PASSWORD_REGEX,
  PATH_LOGIN,
  REGISTER_TITLE,
} from "../../const/const";
import useTitle from "../../hooks/useTitle";

export default function Register() {
  const [username, setUsername] = useState<string>();
  const [isUsernameValid, setIsUsernameValid] = useState<boolean>(true);
  const [email, setEmail] = useState<string>();
  const [isEmailValid, setIsEmailValid] = useState<boolean>(true);
  const [password, setPassword] = useState<string>();
  const [isPasswordValid, setIsPasswordValid] = useState<boolean>(true);
  const [isPasswordShown, setIsPasswordShown] = useState<boolean>(false);

  useTitle(REGISTER_TITLE);

  const navigate = useNavigate();

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_PATH_REGISTER,
    API_METHOD_POST
  );

  function handleUsernameChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(USERNAME_REGEX);
    setUsername(e.target.value);
    setIsUsernameValid(regex.test(e.target.value));
  }

  function handleEmailChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(EMAIL_REGEX);
    setEmail(e.target.value);
    setIsEmailValid(regex.test(e.target.value));
  }

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

    fetchData({ username: username, email: email, password: password }).then(
      (res) => {
        if (res && res.message) {
          navigate(PATH_LOGIN);
        }
      }
    );
  }

  return (
    <AuthFormContainer>
      <h1 className="text-center text-2xl text-primary-black font-bold">
        Register
      </h1>
      <form onSubmit={handleSubmit}>
        {error && <ErrorCard errors={error} />}
        <div className="h-[70px] flex-col gap-2">
          <Input
            type="text"
            name="username"
            placeholder="Username"
            valid={isUsernameValid}
            onChange={handleUsernameChange}
          />
          {!isUsernameValid && (
            <p className="text-xs text-primary-red pl-1 pt-1">
              Username should be 5-12 alphanumeric characters
            </p>
          )}
        </div>
        <div className="mt-3 h-[70px] flex-col gap-2">
          <Input
            type="text"
            name="email"
            placeholder="Email"
            valid={isEmailValid}
            onChange={handleEmailChange}
          />
          {!isEmailValid && (
            <p className="text-xs text-primary-red pl-1 pt-1">
              Email is invalid
            </p>
          )}
        </div>
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
            className={`text-xs ${
              isPasswordValid ? "text-primary-gray" : "text-primary-red"
            } pl-1`}
          >
            Minimum 8 characters with at least one uppercase, one lowercase, one
            special character, and a number
          </p>
        </div>
        <div className="mt-[20px]">
          <Button
            disabled={
              !username ||
              !email ||
              !password ||
              !isUsernameValid ||
              !isEmailValid ||
              !isPasswordValid ||
              isLoading
            }
          >
            {!isLoading ? "Register" : "Loading..."}
          </Button>
        </div>
      </form>

      <p className="text-sm text-center text-primary-black">
        Already have an account?{" "}
        <Link to={PATH_LOGIN} className="text-primary-blue font-bold">
          Login here
        </Link>
      </p>
    </AuthFormContainer>
  );
}
