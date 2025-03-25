import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Eye, EyeOff } from "lucide-react";

import AuthFormContainer from "../../components/AuthFormContainer";
import Input from "../../components/Input";
import ErrorCard from "../../components/ErrorCard";
import Button from "../../components/Button";

import useAxios from "../../hooks/useAxios";
import { SetAccessToken } from "../../utils/localstorage";
import { Response, Token } from "../../types/response";
import {
  API_METHOD_POST,
  API_PATH_LOGIN,
  EMAIL_REGEX,
  LOGIN_TITLE,
  PATH_BACK,
  PATH_FORGET_PASSWORD,
  PATH_REGISTER,
} from "../../const/const";
import useTitle from "../../hooks/useTitle";
import { useSetLoading } from "../../store/user";

export default function Login() {
  const [email, setEmail] = useState<string>();
  const [isEmailValid, setIsEmailValid] = useState<boolean>(true);
  const [password, setPassword] = useState<string>();
  const [isPasswordShown, setIsPasswordShown] = useState<boolean>(false);

  useTitle(LOGIN_TITLE);

  const navigate = useNavigate();

  const setLoading = useSetLoading();

  const { error, isLoading, fetchData } = useAxios<Response<Token>>(
    API_PATH_LOGIN,
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
        navigate(PATH_BACK);
        setLoading();
        setTimeout(() => {
          window.location.reload();
        }, 400);
      }
    });
  }

  return (
    <AuthFormContainer>
      <h1 className="text-center text-2xl font-bold text-primary-black">
        Login
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
          <Link
            to={PATH_FORGET_PASSWORD}
            className="text-sm text-primary-blue pl-1"
          >
            Forget your password?
          </Link>
        </div>
        <div className="mt-[20px]">
          <Button disabled={!email || !password || !isEmailValid || isLoading}>
            Login
          </Button>
        </div>
      </form>

      <p className="text-sm text-center text-primary-black">
        No account yet?{" "}
        <Link to={PATH_REGISTER} className="text-primary-blue font-bold">
          Start your journey here
        </Link>
      </p>
    </AuthFormContainer>
  );
}
