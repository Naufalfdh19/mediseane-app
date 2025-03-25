import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { createPortal } from "react-dom";

import Input from "../../../components/Input";
import Dropdown from "../../../components/Dropdown";
import Button from "../../../components/Button";
import ConfirmBox from "../../../components/ConfirmBox";
import LoaderBox from "../../../components/LoaderBox";
import ErrorBox from "../../../components/ErrorBox";
import SuccessBox from "../../../components/SuccessBox";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import DropdownHeader, {
  searchByType,
} from "../../../components/DropdownHeader";
import ServerError from "../../ServerError";

import { PharmacyData, Response } from "../../../types/response";
import useAxios from "../../../hooks/useAxios";
import {
  API_ADMIN_PHARMACY,
  API_ADMIN_PHARMACIST,
  API_METHOD_GET,
  PATH_ADMIN_PHARMACIST,
  PATH_BACK,
  PHONE_NUMBER_REGEX,
  API_METHOD_PATCH,
  ADMIN_EDIT_PHARMACIST_TITLE,
} from "../../../const/const";
import {
  GetPharmacyByID,
  ParseIDAndNameFromPharmacy,
} from "../../../utils/pharmacy";
import { Pharmacist, PharmacyListItem } from "../../../types/response";
import TitleWithBackButton from "../../../components/TitleWithBackButton";
import ErrorCard from "../../../components/ErrorCard";
import NotFound from "../../NotFound";
import Toggle from "../../../components/Toggle";
import useTitle from "../../../hooks/useTitle";

const searchByList: searchByType[] = [
  {
    id: "1",
    name: "name",
  },
];

export default function EditPharmacist() {
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState("");

  const [name, setName] = useState<string>("");
  const [sipa, setSipa] = useState<string>("");
  const [phoneNumber, setPhoneNumber] = useState<string>("");
  const [isPhoneNumberValid, setIsPhoneNumberValid] = useState<boolean>(true);
  const [years, setYears] = useState<number>(1);
  const [isYearValid, setIsYearValid] = useState<boolean>(true);
  const [email, setEmail] = useState<string>("");
  const [pharmacy, setPharmacy] = useState<PharmacyListItem>();
  const [isAssigned, setIsAssigned] = useState<boolean>(false);
  const [fullPharmacist, setFullPharmacist] = useState<Pharmacist>();

  const [showUpdateConfirm, setShowUpdateConfirm] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  useTitle(ADMIN_EDIT_PHARMACIST_TITLE);

  const navigate = useNavigate();

  const params = useParams();

  const {
    data: pharmacist,
    error: errPharmacist,
    isLoading: isLoadingPharmacist,
  } = useAxios<Response<Pharmacist>>(
    API_ADMIN_PHARMACIST + "/" + params.id,
    API_METHOD_GET
  );

  const {
    data: pharmacies,
    error: errPharmacy,
    isLoading: isLoadingPharmacy,
    fetchData: fetchPharmacy,
  } = useAxios<Response<PharmacyData>>(
    API_ADMIN_PHARMACY +
      "?" +
      (searchByList[0] && filter ? `${searchByList[0].name}=${filter}` : ""),
    API_METHOD_GET
  );

  const { data, error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_ADMIN_PHARMACIST + "/" + params.id,
    API_METHOD_PATCH
  );

  useEffect(() => {
    if (pharmacist?.data) {
      setFullPharmacist(pharmacist.data);
    }
  }, [pharmacist]);

  useEffect(() => {
    if (fullPharmacist) {
      setName(fullPharmacist.name);
      setEmail(fullPharmacist.email);
      setPhoneNumber(fullPharmacist.phone_number);
      setSipa(fullPharmacist.sipa_number);
      setYears(fullPharmacist.year_of_experience);
      setIsAssigned(!!fullPharmacist.pharmacy_id);
      setSearch(fullPharmacist.pharmacy_name);
      setFilter(fullPharmacist.pharmacy_name);
    }
  }, [fullPharmacist]);

  useEffect(() => {
    fetchPharmacy();
  }, [filter]);

  useEffect(() => {
    if (pharmacies && !pharmacy && fullPharmacist?.id) {
      setPharmacy(
        GetPharmacyByID(fullPharmacist.pharmacy_id?.toString(), pharmacies)
      );
    }
  }, [pharmacies]);

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

  function handleSubmit() {
    fetchData({
      id: params.id ? parseInt(params.id) : undefined,
      pharmacy_id: pharmacy?.id && isAssigned ? pharmacy.id : null,
      phone_number: phoneNumber,
      year_of_experience: years,
    });
    setShowResult(true);
    setShowUpdateConfirm(false);
  }

  if (errPharmacist && errPharmacist[0].field === "not found") {
    return <NotFound />;
  }

  if (errPharmacist && errPharmacist[0].field === "server") {
    return <ServerError />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-7 md:my-16 w-[90%] max-w-[1108px] flex flex-col gap-7 md:gap-10">
        <TitleWithBackButton>Edit Pharmacist</TitleWithBackButton>
        <div className="w-full">
          {errPharmacy && <ErrorCard errors={errPharmacy} />}
          <div className="w-full">
            <div className="w-full flex flex-col gap-8">
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Username
                </label>
                {isLoadingPharmacist ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="name"
                    placeholder="Name"
                    type="text"
                    value={name ? name : ""}
                    readOnly
                  />
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="sipa"
                  className="text-primary-black font-semibold"
                >
                  SIPA number
                </label>
                {isLoadingPharmacist ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="sipa"
                    placeholder="SIPA number"
                    type="text"
                    value={sipa ? sipa : ""}
                    readOnly
                  />
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="years"
                  className="text-primary-black font-semibold"
                >
                  Years of experience
                </label>
                {isLoadingPharmacist ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="years"
                    placeholder="1"
                    type="number"
                    value={years?.toString() ? years?.toString() : ""}
                    valid={isYearValid}
                    onChange={handleYearsOfExperienceChange}
                  />
                )}
                {!isYearValid && (
                  <p className="text-sm text-primary-red pl-1 pt-1">
                    Years of experience should be at least 1 year
                  </p>
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="pharmacy"
                  className="text-primary-black font-semibold"
                >
                  Pharmacy
                </label>
                <div className="w-full flex flex-col md:flex-row gap-6">
                  <div className="md:w-[50%]">
                    {isLoadingPharmacy ? (
                      <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                    ) : (
                      <DropdownHeader
                        search={search}
                        setSearch={setSearch}
                        setFilter={setFilter}
                      />
                    )}
                  </div>
                  <div className="md:w-[50%]">
                    {isLoadingPharmacy ? (
                      <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                    ) : (
                      <Dropdown
                        name="pharmacy"
                        placeholder="Select pharmacy"
                        data={ParseIDAndNameFromPharmacy(pharmacies)}
                        selectedId={pharmacy?.id.toString()}
                        onSelect={(id) => {
                          setPharmacy(GetPharmacyByID(id, pharmacies));
                        }}
                        disabled={
                          errPharmacy !== undefined || isLoadingPharmacy
                        }
                      />
                    )}
                  </div>
                  <Toggle
                    labelPosition="right"
                    checked={isAssigned}
                    setChecked={() => {
                      setIsAssigned((prev) => {
                        return !prev;
                      });
                    }}
                  >
                    <p className="text-primary-black font-semibold">
                      {"Assign"}
                    </p>
                  </Toggle>
                </div>
              </div>
              <div className="relative w-full flex flex-col gap-1.5">
                <label
                  htmlFor="phone_number"
                  className="text-primary-black font-semibold"
                >
                  WhatsApp
                </label>
                {isLoadingPharmacist ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="phone_number"
                    placeholder="08xxxxxxxxxx"
                    valid={isPhoneNumberValid}
                    type="text"
                    value={phoneNumber ? phoneNumber : ""}
                    onChange={handlePhoneNumberChange}
                  />
                )}
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
                  Email
                </label>
                {isLoadingPharmacist ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="email"
                    placeholder="pharmacist@mail.com"
                    type="text"
                    value={email ? email : ""}
                    readOnly
                  />
                )}
              </div>
            </div>
          </div>
        </div>
        <div className="w-full flex justify-end">
          <div className="w-[367px] flex justify-between gap-4">
            <Button
              disabled={!(years && phoneNumber) || isLoading}
              onClick={() => {
                setShowUpdateConfirm(true);
              }}
            >
              Save
            </Button>
            <Button
              submit={false}
              type="ghost-green"
              onClick={() => {
                navigate(PATH_BACK);
              }}
            >
              Cancel
            </Button>
          </div>
        </div>
      </div>
      {showUpdateConfirm &&
        createPortal(
          <ConfirmBox
            onYes={handleSubmit}
            onCancel={() => setShowUpdateConfirm(false)}
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
                  : error && CapitalizeFirstLetter(error[0].detail)}
              </p>
            </ErrorBox>
          ) : (
            <SuccessBox onClose={() => navigate(PATH_ADMIN_PHARMACIST)}>
              Pharmacist {data ? "updated" : "deleted"} successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
