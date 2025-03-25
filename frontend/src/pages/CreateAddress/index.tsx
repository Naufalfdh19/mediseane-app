import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import { MapPin } from "lucide-react";

import Button from "../../components/Button";
import Dropdown from "../../components/Dropdown";
import TextArea from "../../components/TextArea";
import MapElem, { MapCoords } from "../../components/MapElem";
import ErrorCard from "../../components/ErrorCard";
import Input from "../../components/Input";
import ConfirmBox from "../../components/ConfirmBox";
import ErrorBox from "../../components/ErrorBox";
import SuccessBox from "../../components/SuccessBox";
import LoaderBox from "../../components/LoaderBox";
import TitleWithBackButton from "../../components/TitleWithBackButton";

import useAxios from "../../hooks/useAxios";
import {
  GetLocationByID,
  ParseIDAndNameFromLocation,
} from "../../utils/address";
import {
  City,
  CurrLocation,
  District,
  Province,
  Response,
  SubDistrict,
} from "../../types/response";
import {
  API_ADDRESS_CITY,
  API_ADDRESS_CHECK,
  API_ADDRESS_DISTRICT,
  API_ADDRESS_PROVINCE,
  API_ADDRESS_SUBDISTRICT,
  API_ADDRESS_USER,
  API_METHOD_GET,
  API_METHOD_POST,
  PHONE_NUMBER_REGEX,
  PATH_BACK,
  CREATE_ADDRESS_TITLE,
} from "../../const/const";
import { CapitalizeFirstLetter } from "../../utils/formatter";
import useTitle from "../../hooks/useTitle";
import ServerError from "../ServerError";

export default function CreateAddress() {
  const [name, setName] = useState<string>();
  const [phoneNumber, setPhoneNumber] = useState<string>();
  const [isPhoneNumberValid, setIsPhoneNumberValid] = useState<boolean>(true);
  const [province, setProvince] = useState<Province>();
  const [city, setCity] = useState<City>();
  const [district, setDistrict] = useState<District>();
  const [subDistrict, setSubDistrict] = useState<SubDistrict>();
  const [addrDetails, setAddrDetails] = useState<string>();
  const [postalCode, setPostalCode] = useState<string>();
  const [center, setCenter] = useState<MapCoords>();
  const [fetchingLoc, setFetchingLoc] = useState<boolean>(false);
  const [locData, setLocData] = useState<CurrLocation>();
  const [showSaveModal, setShowSaveModal] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  useTitle(CREATE_ADDRESS_TITLE);

  const navigate = useNavigate();

  const {
    data: provinces,
    error: errProv,
    isLoading: isLoadingProv,
  } = useAxios<Response<Province[]>>(API_ADDRESS_PROVINCE, API_METHOD_GET);

  const {
    data: cities,
    error: errCity,
    isLoading: isLoadingCity,
    fetchData: fetchCity,
  } = useAxios<Response<City[]>>(
    API_ADDRESS_CITY + province?.id,
    API_METHOD_GET,
    false,
    true
  );

  const {
    data: districts,
    error: errDis,
    isLoading: isLoadingDis,
    fetchData: fetchDistrict,
  } = useAxios<Response<District[]>>(
    API_ADDRESS_DISTRICT + city?.id,
    API_METHOD_GET,
    false,
    true
  );

  const {
    data: subDistricts,
    error: errSubDis,
    isLoading: isLoadingSubDis,
    fetchData: fetchSubDistrict,
  } = useAxios<Response<SubDistrict[]>>(
    API_ADDRESS_SUBDISTRICT + district?.id,
    API_METHOD_GET,
    false,
    true
  );

  const {
    error: errCurrLoc,
    isLoading: isLoadingCurrLoc,
    fetchData: fetchCurrLoc,
  } = useAxios<Response<CurrLocation>>(
    API_ADDRESS_CHECK + `?lat=${center?.lat}&long=${center?.lng}`,
    API_METHOD_GET,
    false,
    true
  );

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_ADDRESS_USER,
    API_METHOD_POST
  );

  useEffect(() => {
    setCity(undefined);
    setCenter(undefined);
    if (province?.id) {
      fetchCity();
    }
  }, [province]);

  useEffect(() => {
    if (cities !== undefined && locData) {
      if (locData.city_id) {
        setCity(GetLocationByID(locData.city_id, cities));
        setLocData(undefined);
      }
    }
  }, [cities]);

  useEffect(() => {
    setDistrict(undefined);
    setCenter(undefined);
    if (city?.id && city.longitude && city.latitude) {
      setCenter({
        lat: parseFloat(city.latitude),
        lng: parseFloat(city.longitude),
      });
      fetchDistrict();
    }
  }, [city]);

  useEffect(() => {
    setSubDistrict(undefined);
    setCenter(undefined);
    if (district?.id) {
      if (district?.latitude && district.longitude) {
        setCenter({
          lat: parseFloat(district.latitude),
          lng: parseFloat(district.longitude),
        });
      }
      fetchSubDistrict();
    }
  }, [district]);

  useEffect(() => {
    setPostalCode("");
  }, [subDistrict]);

  useEffect(() => {
    if (center?.lat && center.lng && fetchingLoc) {
      fetchCurrLoc().then((res) => {
        if (res && res.data) {
          if (res.data.province_id) {
            setLocData(res.data);
            setProvince(GetLocationByID(res.data.province_id, provinces));
          }
        }
      });
    }
    setFetchingLoc(false);
  }, [fetchingLoc]);

  function getCurrentLocation() {
    navigator.geolocation.getCurrentPosition((position) => {
      const lat = position.coords.latitude;
      const lng = position.coords.longitude;
      setCenter({ lat: lat, lng: lng });
      setFetchingLoc(true);
    });
  }

  function handlePhoneNumberChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(PHONE_NUMBER_REGEX);
    setPhoneNumber(e.target.value);
    setIsPhoneNumberValid(regex.test(e.target.value));
  }

  function handleSubmit() {
    fetchData({
      name: name,
      phone_number: phoneNumber,
      address: addrDetails,
      province_id: province?.id,
      province: province?.name,
      city_id: city?.id,
      city: city?.name,
      district_id: district?.id,
      district: district?.name,
      sub_district_id: subDistrict?.id,
      sub_district: subDistrict?.name,
      postal_code: postalCode,
      longitude: center?.lng.toString(),
      latitude: center?.lat.toString(),
    });
    setShowResult(true);
    setShowSaveModal(false);
  }

  if (errProv && errProv[0].field === "server") {
    return <ServerError />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <form className="my-7 md:my-16 w-[90%] max-w-[627px] flex flex-col gap-7 md:gap-10">
        <TitleWithBackButton>Create Address</TitleWithBackButton>
        <div className="w-full flex flex-col gap-5">
          <div className="w-full flex flex-col">
            {(errCurrLoc || errCity || errDis || errSubDis) && (
              <ErrorCard
                errors={
                  errCurrLoc
                    ? errCurrLoc
                    : errCity
                    ? errCity
                    : errDis
                    ? errDis
                    : errSubDis
                    ? errSubDis
                    : []
                }
              />
            )}
            <div className="w-full md:w-[50%] self-end">
              <Button
                onClick={getCurrentLocation}
                disabled={isLoadingCurrLoc}
                submit={false}
              >
                <div className="flex justify-center gap-1.5">
                  <MapPin />
                  <p>Use current location</p>
                </div>
              </Button>
            </div>
          </div>
          <div className="w-full flex flex-col gap-8">
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="name"
                className="text-primary-black font-semibold"
              >
                Name <span className="text-primary-red">*</span>
              </label>
              <Input
                name="name"
                placeholder="Name"
                type="text"
                onChange={(e) => {
                  setName(e.target.value);
                }}
              />
            </div>
            <div className="w-full flex flex-col gap-1.5 relative">
              <label
                htmlFor="phone-number"
                className="text-primary-black font-semibold"
              >
                Phone number <span className="text-primary-red">*</span>
              </label>
              <Input
                name="phone-number"
                placeholder="Phone number"
                type="text"
                valid={isPhoneNumberValid}
                onChange={handlePhoneNumberChange}
              />
              {!isPhoneNumberValid && (
                <p className="absolute bottom-[-22px] text-sm text-primary-red pl-1">
                  Phone number is invalid
                </p>
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="province"
                className="text-primary-black font-semibold"
              >
                Province <span className="text-primary-red">*</span>
              </label>
              <Dropdown
                name="province"
                placeholder="Select province"
                data={ParseIDAndNameFromLocation(provinces)}
                selectedId={province?.id}
                onSelect={(id) => {
                  setProvince(GetLocationByID(id, provinces));
                }}
                disabled={errProv !== undefined || isLoadingProv}
              />
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="city"
                className="text-primary-black font-semibold"
              >
                City <span className="text-primary-red">*</span>
              </label>
              <Dropdown
                name="city"
                placeholder="Select city"
                data={ParseIDAndNameFromLocation(cities)}
                selectedId={city?.id}
                onSelect={(id) => {
                  setCity(GetLocationByID(id, cities));
                }}
                disabled={!province || errCity !== undefined || isLoadingCity}
              />
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="district"
                className="text-primary-black font-semibold"
              >
                District <span className="text-primary-red">*</span>
              </label>
              <Dropdown
                name="district"
                placeholder="Select district"
                data={ParseIDAndNameFromLocation(districts)}
                hasImage={false}
                selectedId={district?.id}
                onSelect={(id) => {
                  setDistrict(GetLocationByID(id, districts));
                }}
                disabled={!city || errDis !== undefined || isLoadingDis}
              />
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="subDistrict"
                className="text-primary-black font-semibold"
              >
                Sub district <span className="text-primary-red">*</span>
              </label>
              <Dropdown
                name="subDistrict"
                placeholder="Select sub district"
                data={ParseIDAndNameFromLocation(subDistricts)}
                hasImage={false}
                selectedId={subDistrict?.id}
                onSelect={(id) => {
                  setSubDistrict(GetLocationByID(id, subDistricts));
                }}
                disabled={
                  !district || errSubDis !== undefined || isLoadingSubDis
                }
              />
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="postCode"
                className="text-primary-black font-semibold"
              >
                Postal code <span className="text-primary-red">*</span>
              </label>
              <Dropdown
                name="postcode"
                placeholder="Select postal code"
                data={
                  subDistrict?.postal_codes
                    ? subDistrict.postal_codes.split(",").map((item) => {
                        return {
                          id: item,
                          name: item,
                        };
                      })
                    : []
                }
                hasImage={false}
                selectedId={postalCode}
                onSelect={(id) => {
                  setPostalCode(id);
                }}
                disabled={!subDistrict}
              />
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="address-details"
                className="text-primary-black font-semibold"
              >
                Address details <span className="text-primary-red">*</span>
              </label>
              <TextArea
                name="address-details"
                placeholder="Street name, building number, etc."
                onChange={(e) => {
                  setAddrDetails(e.target.value);
                }}
              />
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <p className="text-primary-black font-semibold">
                Map pinpoint <span className="text-primary-red">*</span>
              </p>
              <MapElem
                center={center}
                setCenter={setCenter}
                disabled={!center}
              />
            </div>
          </div>
        </div>
        <div className="flex justify-center">
          <div className="w-full max-w-[367px]">
            <Button
              disabled={
                !(
                  name &&
                  phoneNumber &&
                  isPhoneNumberValid &&
                  province &&
                  city &&
                  district &&
                  subDistrict &&
                  addrDetails &&
                  postalCode
                ) || isLoading
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
            <SuccessBox onClose={() => navigate(PATH_BACK)}>
              Address added successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
