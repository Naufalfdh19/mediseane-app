import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import { Eye, EyeOff } from "lucide-react";

import Input from "../../../components/Input";
import Button from "../../../components/Button";
import ConfirmBox from "../../../components/ConfirmBox";
import LoaderBox from "../../../components/LoaderBox";
import ErrorBox from "../../../components/ErrorBox";
import SuccessBox from "../../../components/SuccessBox";
import ErrorCard from "../../../components/ErrorCard";
import TitleWithBackButton from "../../../components/TitleWithBackButton";

import {
  API_ADMIN_PHARMACIST,
  API_METHOD_GET,
  API_METHOD_POST,
  PATH_ADMIN_PHARMACIST,
  PHONE_NUMBER_REGEX,
  API_ADMIN_GENERATED_PASSWORD,
  USERNAME_REGEX,
  ADMIN_ADD_PHARMACIST_TITLE,
} from "../../../const/const";
import useAxios from "../../../hooks/useAxios";
import { GeneratedPassword, Response } from "../../../types/response";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import useTitle from "../../../hooks/useTitle";

export default function AddPharmacist() {
  const [name, setName] = useState<string>();
  const [isUsernameValid, setIsUsernameValid] = useState<boolean>(true);
  const [sipa, setSipa] = useState<string>();
  const [isSipaValid, setIsSipaValid] = useState<boolean>(true);
  const [phoneNumber, setPhoneNumber] = useState<string>();
  const [isPhoneNumberValid, setIsPhoneNumberValid] = useState<boolean>(true);
  const [years, setYears] = useState<number>();
  const [isYearValid, setIsYearValid] = useState<boolean>(true);
  const [email, setEmail] = useState<string>();
  const [password, setPassword] = useState<string>();
  const [isPasswordShown, setIsPasswordShown] = useState<boolean>(false);

  const [showSaveModal, setShowSaveModal] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  useTitle(ADMIN_ADD_PHARMACIST_TITLE)

  const navigate = useNavigate();

  const {
    data: generatedPassword,
    error: errPassword,
    isLoading: isLoadingPassword,
  } = useAxios<Response<GeneratedPassword>>(
    API_ADMIN_GENERATED_PASSWORD,
    API_METHOD_GET
  );

  const { error, isLoading, fetchData } = useAxios(
    API_ADMIN_PHARMACIST,
    API_METHOD_POST
  );

  useEffect(() => {
    if (generatedPassword?.data) {
      setPassword(generatedPassword.data.password);
    }
  }, [generatedPassword]);

  function handleSubmit() {
    fetchData({
      name: name,
      sipa_number: sipa,
      phone_number: phoneNumber,
      year_of_experience: years,
      email: email,
      password: password,
    });
    setShowResult(true);
    setShowSaveModal(false);
  }

  function handleNameChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(USERNAME_REGEX);
    setName(e.target.value);
    setIsUsernameValid(regex.test(e.target.value));
  }

  function handleEmailChange(e: React.ChangeEvent<HTMLInputElement>) {
    setEmail(e.target.value);
  }

  function handleSipaChange(e: React.ChangeEvent<HTMLInputElement>) {
    const SIPA_NUMBER_MINLENGTH = 9;
    const sipaNumber = e.target.value;
    setSipa(sipaNumber);
    setIsSipaValid(sipaNumber.length >= SIPA_NUMBER_MINLENGTH);
  }

  function handleYearsOfExperienceChange(
    e: React.ChangeEvent<HTMLInputElement>
  ) {
    const years = Number(e.target.value);
    setYears(years);
    setIsYearValid(years > 0);
  }

  function handlePhoneNumberChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(PHONE_NUMBER_REGEX);
    setPhoneNumber(e.target.value);
    setIsPhoneNumberValid(regex.test(e.target.value));
  }

  function handlePasswordToggleClick() {
    setIsPasswordShown((prev) => !prev);
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full">
          <TitleWithBackButton>Add Pharmacist</TitleWithBackButton>
        </div>
        <div className="w-full">
          <div className="w-full flex flex-col gap-8">
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="name"
                className="text-primary-black font-semibold"
              >
                Username <span className="text-primary-red">*</span>
              </label>
              <Input
                name="name"
                placeholder="Username"
                type="text"
                valid={isUsernameValid}
                onChange={handleNameChange}
              />
              {!isUsernameValid && (
                <p className="text-sm text-primary-red pl-1 pt-1">
                  Username should be 5-12 alphanumeric characters
                </p>
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="sipa"
                className="text-primary-black font-semibold"
              >
                SIPA number <span className="text-primary-red">*</span>
              </label>
              <Input
                name="sipa"
                placeholder="SIPA number"
                type="text"
                valid={isSipaValid}
                onChange={handleSipaChange}
              />
              {!isSipaValid && (
                <p className="text-sm text-primary-red pl-1 pt-1">
                  SIPA number should be at least 9 characters
                </p>
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="years"
                className="text-primary-black font-semibold"
              >
                Years of experience <span className="text-primary-red">*</span>
              </label>
              <Input
                name="years"
                placeholder="1"
                type="number"
                valid={isYearValid}
                onChange={handleYearsOfExperienceChange}
              />
              {!isYearValid && (
                <p className="text-sm text-primary-red pl-1 pt-1">
                  Years of experience should be at least 1 year
                </p>
              )}
            </div>
            <div className="relative w-full flex flex-col gap-1.5">
              <label
                htmlFor="phone_number"
                className="text-primary-black font-semibold"
              >
                WhatsApp <span className="text-primary-red">*</span>
              </label>
              <Input
                name="phone_number"
                placeholder="08xxxxxxxxxx"
                valid={isPhoneNumberValid}
                type="text"
                onChange={handlePhoneNumberChange}
              />
              {!isPhoneNumberValid && (
                <p className="absolute bottom-[-22px] text-sm text-primary-red pl-1">
                  WhatsApp number is invalid
                </p>
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="email"
                className="text-primary-black font-semibold"
              >
                Email <span className="text-primary-red">*</span>
              </label>
              <Input
                name="email"
                placeholder="pharmacist@mail.com"
                type="text"
                onChange={handleEmailChange}
              />
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="password"
                className="text-primary-black font-semibold"
              >
                Auto-generated password{" "}
                <span className="text-primary-red">*</span>
              </label>
              {isLoadingPassword ? (
                <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
              ) : errPassword || !password ? (
                <ErrorCard errors={errPassword ? errPassword : []} />
              ) : (
                <Input
                  type={isPasswordShown ? "text" : "password"}
                  name="password"
                  placeholder={password}
                  value={password}
                  icon={isPasswordShown ? <Eye /> : <EyeOff />}
                  onIconClick={handlePasswordToggleClick}
                  readOnly
                />
              )}
            </div>
          </div>
        </div>
        <div className="w-full flex justify-end">
          <div className="w-[367px] flex justify-between gap-4">
            <Button
              disabled={
                !(
                  name &&
                  sipa &&
                  phoneNumber &&
                  isPhoneNumberValid &&
                  years &&
                  email
                ) || isLoading
              }
              onClick={() => {
                setShowSaveModal(true);
              }}
            >
              Save
            </Button>
            <Button
              submit={false}
              type="ghost-green"
              onClick={() => {
                navigate(PATH_ADMIN_PHARMACIST);
              }}
            >
              Cancel
            </Button>
          </div>
        </div>
      </div>
      {showSaveModal &&
        createPortal(
          <ConfirmBox
            onYes={handleSubmit}
            onCancel={() => setShowSaveModal(false)}
          />,
          document.body
        )}
      {showResult &&
        createPortal(
          isLoading ? (
            <LoaderBox />
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
            <SuccessBox onClose={() => navigate(PATH_ADMIN_PHARMACIST)}>
              Pharmacist added successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
