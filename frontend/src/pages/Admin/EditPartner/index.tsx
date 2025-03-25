import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { useNavigate, useParams } from "react-router-dom";
import { Trash2 } from "lucide-react";

import TitleWithBackButton from "../../../components/TitleWithBackButton";
import Input from "../../../components/Input";
import Dropdown from "../../../components/Dropdown";
import Toggle from "../../../components/Toggle";
import Button from "../../../components/Button";
import ConfirmBox from "../../../components/ConfirmBox";
import LoaderBox from "../../../components/LoaderBox";
import ErrorBox from "../../../components/ErrorBox";
import SuccessBox from "../../../components/SuccessBox";
import PageLoader from "../../../components/PageLoader";
import ServerError from "../../ServerError";

import useAxios from "../../../hooks/useAxios";
import { IsHourRangeValid } from "../../../utils/partners";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import { Partner, Response } from "../../../types/response";
import {
  ADMIN_EDIT_PARTNER_TITLE,
  API_ADMIN_PARTNERS,
  API_METHOD_DELETE,
  API_METHOD_PATCH,
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

export default function EditPartner() {
  const [name, setName] = useState<string>();
  const [year, setYear] = useState<string>();
  const [activeDays, setActiveDays] = useState<number[]>([]);
  const [start, setStart] = useState<string>();
  const [end, setEnd] = useState<string>();
  const [isOperationalHoursValid, setIsOperationalHoursValid] =
    useState<boolean>(true);
  const [isActive, setIsActive] = useState<boolean>(false);
  const [showUpdateConfirm, setShowUpdateConfirm] = useState<boolean>(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  useTitle(ADMIN_EDIT_PARTNER_TITLE);

  const params = useParams();

  const { data, error, isLoading } = useAxios<Response<Partner>>(
    API_ADMIN_PARTNERS + "/" + params.id
  );

  const {
    data: updateData,
    error: errorUpdate,
    isLoading: isUpdateLoading,
    fetchData,
  } = useAxios(API_ADMIN_PARTNERS + "/" + params.id, API_METHOD_PATCH);

  const {
    error: errorDelete,
    isLoading: isDeleteLoading,
    fetchData: fetchDelete,
  } = useAxios(API_ADMIN_PARTNERS + "/" + params.id, API_METHOD_DELETE);

  const navigate = useNavigate();

  useEffect(() => {
    if (data) {
      setName(data.data.name);
      setYear(data.data.year_founded);
      setStart(data.data.start_hour);
      setEnd(data.data.end_hour);
      setIsActive(data.data.is_active);
      const activeDay: number[] = [];
      const dayArr = data.data.active_days.split(",");
      dayArr.forEach((day) => {
        activeDay.push(DAYS.findIndex((item) => item === day));
      });
      setActiveDays(activeDay.sort());
    }
  }, [data]);

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
    setShowUpdateConfirm(false);
    setShowResult(true);
  }

  function handleDelete() {
    fetchDelete();
    setShowResult(true);
  }

  if (isLoading) {
    return <PageLoader />;
  }

  if (error) {
    return <ServerError />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[800px] flex flex-col gap-8">
        <TitleWithBackButton>Edit Partner</TitleWithBackButton>
        <div className="grow">
          <form className="flex flex-col gap-8">
            <div className="flex flex-col md:flex-row gap-8">
              <div className="flex-1 flex flex-col gap-1.5">
                <label
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Name
                </label>
                <Input
                  type="text"
                  name="name"
                  placeholder="name"
                  value={name ? name : ""}
                  readOnly={true}
                  onChange={handleNameChange}
                />
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
                  disabled={true}
                />
              </div>
            </div>

            <div className="flex-1 flex flex-col gap-2 relative">
              <p className="text-primary-black font-semibold">Active days</p>
              <div className="flex flex-wrap gap-2">
                {DAYS.map((day, idx) => {
                  return (
                    <div className="flex" key={day}>
                      <input
                        key={idx}
                        type="checkbox"
                        name={day}
                        id={day}
                        value={idx}
                        checked={activeDays.includes(idx)}
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
                      value={start ? start : ""}
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
                      value={end ? end : ""}
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

            <div className="flex gap-4 justify-center">
              <Button
                disabled={
                  !(
                    name &&
                    year &&
                    activeDays.length > 0 &&
                    start &&
                    end &&
                    isOperationalHoursValid
                  ) || isLoading
                }
                submit={false}
                onClick={() => {
                  setShowUpdateConfirm(true);
                }}
              >
                Save
              </Button>
              <Button
                type="ghost-red"
                onClick={() => {
                  setShowDeleteConfirm(true);
                }}
                submit={false}
              >
                <div className="flex justify-center gap-1.5 text-primary-red">
                  <Trash2 />
                  <p>Delete</p>
                </div>
              </Button>
            </div>
          </form>
          {showUpdateConfirm &&
            createPortal(
              <ConfirmBox
                type="update"
                onYes={handleSubmit}
                onCancel={() => setShowUpdateConfirm(false)}
              />,
              document.body
            )}
          {showDeleteConfirm &&
            createPortal(
              <ConfirmBox
                type="delete"
                onYes={() => {
                  handleDelete();
                  setShowDeleteConfirm(false);
                }}
                onCancel={() => {
                  setShowDeleteConfirm(false);
                }}
              />,
              document.body
            )}
          {showResult &&
            createPortal(
              isUpdateLoading || isDeleteLoading ? (
                <LoaderBox />
              ) : errorUpdate || errorDelete ? (
                <ErrorBox onClose={() => setShowResult(false)}>
                  <p className="font-bold text-primary-black text-4xl">
                    Oops...
                  </p>
                  <p className="text-primary-black text-xl">
                    {(errorUpdate && errorUpdate[0].field === "server") ||
                    (errorDelete && errorDelete[0].field === "server")
                      ? "Something is wrong, please try again"
                      : (errorUpdate &&
                          CapitalizeFirstLetter(errorUpdate[0].detail)) ||
                        (errorDelete &&
                          CapitalizeFirstLetter(errorDelete[0].detail))}
                  </p>
                </ErrorBox>
              ) : (
                <SuccessBox onClose={() => navigate(PATH_BACK)}>
                  Partener {updateData ? "updated" : "deleted"} successfully!
                </SuccessBox>
              ),
              document.body
            )}
        </div>
      </div>
    </div>
  );
}
