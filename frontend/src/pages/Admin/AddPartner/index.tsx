import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { useNavigate } from "react-router-dom";

import TitleWithBackButton from "../../../components/TitleWithBackButton";
import Input from "../../../components/Input";
import Dropdown from "../../../components/Dropdown";
import Toggle from "../../../components/Toggle";
import Button from "../../../components/Button";
import ConfirmBox from "../../../components/ConfirmBox";
import LoaderBox from "../../../components/LoaderBox";
import ErrorBox from "../../../components/ErrorBox";
import SuccessBox from "../../../components/SuccessBox";

import useAxios from "../../../hooks/useAxios";
import { IsHourRangeValid } from "../../../utils/partners";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import {
  ADMIN_ADD_PARTNER_TITLE,
  API_ADMIN_PARTNERS,
  API_METHOD_POST,
  DAYS,
  PATH_BACK,
} from "../../../const/const";
import useTitle from "../../../hooks/useTitle";

interface dropdownItem {
  id: string;
  name: string;
}

const currentYear = new Date().getFullYear();
const years: dropdownItem[] = [];
for (let year = currentYear - 100; year <= currentYear; year++) {
  years.push({ id: year.toString(), name: year.toString() });
}

export default function AddPartner() {
  const [name, setName] = useState<string>();
  const [isNameValid, setIsNameValid] = useState<boolean>(true);
  const [year, setYear] = useState<string>();
  const [activeDays, setActiveDays] = useState<number[]>([]);
  const [start, setStart] = useState<string>();
  const [end, setEnd] = useState<string>();
  const [isOperationalHoursValid, setIsOperationalHoursValid] =
    useState<boolean>(true);
  const [isActive, setIsActive] = useState<boolean>(false);
  const [showSaveModal, setShowSaveModal] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  useTitle(ADMIN_ADD_PARTNER_TITLE);

  const { error, isLoading, fetchData } = useAxios(
    API_ADMIN_PARTNERS,
    API_METHOD_POST
  );

  const navigate = useNavigate();

  useEffect(() => {
    if (start && end) {
      setIsOperationalHoursValid(IsHourRangeValid(start, end));
    }
  }, [start, end]);

  function handleDaysChange(e: React.ChangeEvent<HTMLInputElement>) {
    setActiveDays((prev) => {
      if (prev.includes(+e.target.value)) {
        const newArr = activeDays.filter((day) => day !== +e.target.value);
        return newArr;
      } else {
        return [...prev, +e.target.value];
      }
    });
  }

  function handleNameChange(e: React.ChangeEvent<HTMLInputElement>) {
    setName(e.target.value);
    setIsNameValid(e.target.value.length > 2);
  }

  function handleStartChange(e: React.ChangeEvent<HTMLInputElement>) {
    setStart(e.target.value);
  }
  function handleEndChange(e: React.ChangeEvent<HTMLInputElement>) {
    setEnd(e.target.value);
  }

  function handleSubmit() {
    let activeDay: string = "";
    const sortedDays = activeDays.sort();
    sortedDays.forEach((day, idx) => {
      if (idx === sortedDays.length - 1) {
        activeDay += DAYS[day];
      } else {
        activeDay += `${DAYS[day]},`;
      }
    });
    fetchData({
      name: name,
      year_founded: year,
      active_days: activeDay,
      start_hour: start,
      end_hour: end,
      is_active: isActive,
    });
    setShowSaveModal(false);
    setShowResult(true);
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[800px] flex flex-col gap-8">
        <TitleWithBackButton>Add Partner</TitleWithBackButton>
        <div className="grow">
          <form className="flex flex-col gap-8">
            <div className="flex flex-col md:flex-row gap-8">
              <div className="flex-1 flex flex-col gap-1.5 relative">
                <label
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Name
                </label>
                <Input
                  type="text"
                  name="name"
                  placeholder="Name"
                  onChange={handleNameChange}
                  valid={isNameValid}
                />
                {!isNameValid && (
                  <p className="absolute bottom-[-25px] text-primary-red">
                    Name should have at least 3 characters
                  </p>
                )}
              </div>
              <div className="w-[130px] flex flex-col gap-1.5">
                <label
                  htmlFor="year"
                  className="text-primary-black font-semibold"
                >
                  Year founded
                </label>
                <Dropdown
                  name="year"
                  placeholder="YYYY"
                  data={years}
                  selectedId={year}
                  onSelect={(year) => setYear(year)}
                />
              </div>
            </div>

            <div className="flex-1 relative">
              <p className="text-primary-black font-semibold">Active days</p>
              <div className="flex flex-wrap gap-2 mt-2">
                {DAYS.map((day, idx) => {
                  return (
                    <div className="flex" key={day}>
                      <input
                        key={idx}
                        type="checkbox"
                        name={day}
                        id={day}
                        value={idx}
                        className="mr-2 accent-primary-green"
                        onChange={handleDaysChange}
                      />
                      <label htmlFor={day}>{CapitalizeFirstLetter(day)}</label>
                    </div>
                  );
                })}
              </div>
            </div>

            <div className="flex flex-col md:flex-row gap-8 justify-between">
              <div className="flex-1 relative">
                <p className="text-primary-black font-semibold">
                  Operational hours
                </p>
                <div className="flex flex-col md:flex-row gap-4">
                  <div className="flex-1 flex flex-col gap-1.5">
                    <label
                      htmlFor="from"
                      className="text-primary-black font-semibold"
                    >
                      Start
                    </label>
                    <Input
                      type="time"
                      name="start"
                      placeholder="start"
                      valid={isOperationalHoursValid}
                      onChange={handleStartChange}
                    />
                  </div>
                  <div className="flex-1 flex flex-col gap-1.5">
                    <label
                      htmlFor="to"
                      className="text-primary-black font-semibold"
                    >
                      End
                    </label>
                    <Input
                      type="time"
                      name="end"
                      placeholder="end"
                      valid={isOperationalHoursValid}
                      onChange={handleEndChange}
                    />
                  </div>
                </div>
                {!isOperationalHoursValid && (
                  <p className="absolute bottom-[-26px] text-primary-red">
                    Invalid operational hours
                  </p>
                )}
              </div>
            </div>

            <Toggle
              labelPosition="right"
              checked={isActive}
              setChecked={() => setIsActive((prev) => !prev)}
            >
              <p className="text-primary-black font-semibold">Active</p>
            </Toggle>
            <div className="flex justify-center">
              <div className="w-full max-w-[367px]">
                <Button
                  disabled={
                    !(
                      name &&
                      isNameValid &&
                      year &&
                      activeDays.length > 0 &&
                      start &&
                      end &&
                      isOperationalHoursValid
                    )
                  }
                  submit={false}
                  onClick={() => {
                    setShowSaveModal(true);
                  }}
                >
                  Save
                </Button>
              </div>
            </div>
          </form>
          {showSaveModal &&
            createPortal(
              <ConfirmBox
                type="save"
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
                  <p className="font-bold text-primary-black text-4xl">
                    Oops...
                  </p>
                  <p className="text-primary-black text-xl">
                    {error && error[0].field === "server"
                      ? "Something is wrong, please try again"
                      : CapitalizeFirstLetter(error[0].detail)}
                  </p>
                </ErrorBox>
              ) : (
                <SuccessBox onClose={() => navigate(PATH_BACK)}>
                  Partner added successfully!
                </SuccessBox>
              ),
              document.body
            )}
        </div>
      </div>
    </div>
  );
}
