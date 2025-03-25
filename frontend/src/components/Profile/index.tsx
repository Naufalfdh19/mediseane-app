import { useState } from "react";
import Input from "../Input";
import {
  API_METHOD_POST,
  API_PATH_REGISTER,
  EMAIL_REGEX,
  PASSWORD_REGEX,
  PATH_ADDRESS,
  PATH_HOME,
  PATH_LOGIN,
  USERNAME_REGEX,
} from "../../const/const";
import useAxios from "../../hooks/useAxios";
import { Response } from "../../types/response";
import ErrorCard from "../ErrorCard";
import {
  ChevronRight,
  CircleUserRound,
  Eye,
  EyeOff,
  LogOut,
} from "lucide-react";
import Button from "../Button";
import { Link, useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import ConfirmBox from "../ConfirmBox";
import { RemoveAccessToken } from "../../utils/localstorage";

export default function Profile() {
  const [username, setUsername] = useState<string>();
  const [isUsernameValid, setIsUsernameValid] = useState<boolean>(true);
  const [email, setEmail] = useState<string>();
  const [isEmailValid, setIsEmailValid] = useState<boolean>(true);
  const [password, setPassword] = useState<string>();
  const [isPasswordValid, setIsPasswordValid] = useState<boolean>(true);
  const [isPasswordShown, setIsPasswordShown] = useState<boolean>(false);
  const [showLogoutConfirm, setShowLogoutConfirm] = useState<boolean>(false);

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
        }
      }
    );
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[627px] flex flex-col items-center gap-10">
        <div className="w-full flex flex-col gap-8 items-center">
          <h1 className="text-center text-2xl font-bold text-primary-black">
            Profile
          </h1>
          <form
            onSubmit={handleSubmit}
            className="flex flex-col items-center mb-4"
          >
            <CircleUserRound className="w-full max-w-[196px] h-full max-h-[196px] mb-4 text-secondary-gray" />
            {error && <ErrorCard errors={error} />}
            <div className="flex flex-col gap-3">
              <div className="flex flex-col gap-2">
                <h1 className="text-primary-black font-semibold">Username</h1>
                <Input
                  type="text"
                  name="username"
                  placeholder="Username"
                  valid={isUsernameValid}
                  onChange={handleUsernameChange}
                />
                {!isUsernameValid && (
                  <p className="text-sm text-primary-red pl-1 pt-1">
                    Username should be 5-12 alphanumeric characters
                  </p>
                )}
              </div>
              <div className=" flex flex-col gap-2">
                <h1 className="text-primary-black font-semibold">Email</h1>
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
                <p className={"text-sm text-primary-gray pl-1"}>
                  Order product?{" "}
                  <Link to="" className="text-primary-blue font-bold">
                    Let’s verify email first!
                  </Link>
                </p>
              </div>
              <div className="flex flex-col gap-2">
                <h1 className="text-primary-black font-semibold">Password</h1>
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
                  Minimum 8 characters with at least one uppercase, one
                  lowercase, one special character, and a number
                </p>
              </div>
            </div>
            <div className="w-full mt-8 flex justify-center">
              <div className="w-full max-w-[367px] flex justify-center gap-4">
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
                  Save
                </Button>
                <Button
                  type="ghost"
                  onClick={() => {
                    navigate(PATH_HOME);
                  }}
                  submit={false}
                >
                  Cancel
                </Button>
              </div>
            </div>
          </form>
          <div className="w-full flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <h1 className="text-primary-black font-semibold">Address</h1>
              <Button
                onClick={() => {
                  navigate(PATH_ADDRESS);
                }}
                type="ghost"
                size="xl"
                square={true}
              >
                <div className="px-3 flex justify-between">
                  <h1 className="text-primary-blue font-bold">
                    See saved addresses
                  </h1>
                  <ChevronRight className="text-primary-blue" />
                </div>
              </Button>
            </div>
            <div className="flex flex-col gap-2">
              <h1 className="text-primary-black font-semibold">Logout</h1>
              <Button
                onClick={() => {
                  setShowLogoutConfirm(true);
                }}
                type="ghost"
                square={true}
              >
                <div className="px-3 flex justify-between">
                  <h1 className="text-primary-red font-bold">Logout</h1>
                  <LogOut className="text-primary-red" />
                </div>
              </Button>
            </div>
          </div>
        </div>
      </div>
      {showLogoutConfirm &&
        createPortal(
          <ConfirmBox
            type="logout"
            onYes={() => {
              RemoveAccessToken();
              navigate(PATH_LOGIN);
            }}
            onCancel={() => {
              setShowLogoutConfirm(false);
            }}
          />,
          document.body
        )}
    </div>
  );
}
