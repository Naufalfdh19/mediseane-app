import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { createPortal } from "react-dom";
import { Trash2 } from "lucide-react";

import Button from "../../components/Button";
import Dropdown from "../../components/Dropdown";
import TextArea from "../../components/TextArea";
import MapElem, { MapCoords } from "../../components/MapElem";
import ErrorCard from "../../components/ErrorCard";
import Input from "../../components/Input";
import SuccessBox from "../../components/SuccessBox";
import Toggle from "../../components/Toggle";
import ConfirmBox from "../../components/ConfirmBox";
import NotFound from "../NotFound";
import PageLoader from "../../components/PageLoader";
import LoaderBox from "../../components/LoaderBox";
import ErrorBox from "../../components/ErrorBox";
import TitleWithBackButton from "../../components/TitleWithBackButton";

import useAxios from "../../hooks/useAxios";
import { CapitalizeFirstLetter } from "../../utils/formatter";
import {
  GetLocationByID,
  ParseIDAndNameFromLocation,
} from "../../utils/address";
import {
  City,
  District,
  UserAddressWithID,
  Province,
  Response,
  SubDistrict,
} from "../../types/response";
import {
  API_ADDRESS_CITY,
  API_ADDRESS_DISTRICT,
  API_ADDRESS_PROVINCE,
  API_ADDRESS_SUBDISTRICT,
  API_ADDRESS_USER,
  API_METHOD_DELETE,
  API_METHOD_GET,
  API_METHOD_PUT,
  EDIT_ADDRESS_TITLE,
  PATH_BACK,
  PHONE_NUMBER_REGEX,
} from "../../const/const";
import useTitle from "../../hooks/useTitle";
import ServerError from "../ServerError";

export default function EditAddress() {
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
  const [isActive, setIsActive] = useState<boolean>(false);
  const [fullAddress, setFullAddress] = useState<UserAddressWithID>();
  const [showResult, setShowResult] = useState<boolean>(false);
  const [showUpdateConfirm, setShowUpdateConfirm] = useState<boolean>(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const [showNoActiveAddressErr, setShowNoActiveAddressErr] =
    useState<boolean>(false);

  useTitle(EDIT_ADDRESS_TITLE);

  const navigate = useNavigate();

  const params = useParams();

  const {
    data: userData,
    error: errUserData,
    isLoading: isLoadingUserData,
  } = useAxios<Response<UserAddressWithID>>(
    API_ADDRESS_USER + "/" + params.id,
    API_METHOD_GET
  );

  const {
    error: errDelete,
    isLoading: isDeleteLoading,
    fetchData: fetchDelete,
  } = useAxios<Response<undefined>>(
    API_ADDRESS_USER + "/" + params.id,
    API_METHOD_DELETE
  );

  const {
    data: provinces,
    error: errProv,
    isLoading: isLoadingProv,
    fetchData: fetchProvince,
  } = useAxios<Response<Province[]>>(
    API_ADDRESS_PROVINCE,
    API_METHOD_GET,
    false,
    true
  );

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

  const { data, error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_ADDRESS_USER,
    API_METHOD_PUT
  );

  useEffect(() => {
    if (userData?.data) {
      setFullAddress(userData.data);
    }
  }, [userData]);

  useEffect(() => {
    if (fullAddress) {
      setName(fullAddress.name);
      setPhoneNumber(fullAddress.phone_number);
      setAddrDetails(fullAddress.address);
      setIsActive(fullAddress.is_active);
      fetchProvince();
    }
  }, [fullAddress]);

  useEffect(() => {
    if (provinces && !province && fullAddress?.province_id) {
      setProvince(
        GetLocationByID(fullAddress.province_id.toString(), provinces)
      );
    }
  }, [provinces]);

  useEffect(() => {
    if (cities && !city && fullAddress?.city_id) {
      setCity(GetLocationByID(fullAddress.city_id.toString(), cities));
    }
  }, [cities]);

  useEffect(() => {
    if (districts && !district && fullAddress?.district_id) {
      setDistrict(
        GetLocationByID(fullAddress.district_id.toString(), districts)
      );
    }
  }, [districts]);

  useEffect(() => {
    if (subDistricts && !subDistrict && fullAddress?.sub_district_id) {
      setSubDistrict(
        GetLocationByID(fullAddress.sub_district_id.toString(), subDistricts)
      );
    }
  }, [subDistricts]);

  useEffect(() => {
    setCity(undefined);
    setCenter(undefined);
    if (province?.id) {
      fetchCity();
    }
  }, [province]);

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
    if (fullAddress) {
      setPostalCode(fullAddress.postal_code);
      setCenter({
        lat: parseFloat(fullAddress.latitude),
        lng: parseFloat(fullAddress.longitude),
      });
      setFullAddress(undefined);
      setLoading(false);
    }
  }, [subDistrict]);

  function handlePhoneNumberChange(e: React.ChangeEvent<HTMLInputElement>) {
    const regex = new RegExp(PHONE_NUMBER_REGEX);
    setPhoneNumber(e.target.value);
    setIsPhoneNumberValid(regex.test(e.target.value));
  }

  function isDataChanged() {
    let isChanged: boolean = false;
    if (
      userData?.data.name !== name ||
      userData?.data.phone_number !== phoneNumber ||
      userData?.data.province_id !== province?.id ||
      userData?.data.city_id !== city?.id ||
      userData?.data.district_id !== district?.id ||
      userData?.data.sub_district_id !== subDistrict?.id ||
      userData?.data.postal_code !== postalCode ||
      userData?.data.address !== addrDetails ||
      userData?.data.longitude !== center?.lng.toString() ||
      userData?.data.latitude !== center?.lat.toString() ||
      userData?.data.is_active !== isActive
    ) {
      isChanged = true;
    }
    return isChanged;
  }

  function handleDelete() {
    fetchDelete();
    setShowResult(true);
  }

  function handleSubmit() {
    fetchData({
      id: params.id ? parseInt(params.id) : undefined,
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
      is_active: isActive,
    });
    setShowResult(true);
  }

  if (errUserData && errUserData[0].field === "not found") {
    return <NotFound />;
  }

  if (errUserData && errUserData[0].field === "server") {
    return <ServerError />;
  }

  if (isLoadingUserData || loading) {
    return <PageLoader />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <form
        className="my-7 md:my-16 w-[90%] max-w-[627px] flex flex-col gap-7 md:gap-10"
        onSubmit={handleSubmit}
      >
        <TitleWithBackButton>Edit Address</TitleWithBackButton>
        <div className="w-full">
          {(errProv || errCity || errDis || errSubDis) && (
            <ErrorCard
              errors={
                errProv
                  ? errProv
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
          <div className="w-full flex flex-col gap-8">
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="name"
                className="text-primary-black font-semibold"
              >
                Name
              </label>
              <Input
                name="name"
                placeholder="Name"
                type="text"
                value={name}
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
                Phone number
              </label>
              <Input
                name="phone-number"
                placeholder="Phone number"
                type="text"
                value={phoneNumber}
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
                Province
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
                City
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
                District
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
                htmlFor="sub-district"
                className="text-primary-black font-semibold"
              >
                Sub district
              </label>
              <Dropdown
                name="sub-district"
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
                htmlFor="postcode"
                className="text-primary-black font-semibold"
              >
                Postal code
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
                Address details
              </label>
              <TextArea
                name="address-details"
                placeholder="Street Name, Building Number, etc."
                value={addrDetails}
                onChange={(e) => {
                  setAddrDetails(e.target.value);
                }}
              />
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <p className="text-primary-black font-semibold">Map pinpoint</p>
              <MapElem
                center={center}
                setCenter={setCenter}
                disabled={!center}
              />
            </div>
            <div
              className="flex"
              onClick={() => {
                if (userData?.data.is_active) {
                  setShowNoActiveAddressErr(true);
                }
              }}
            >
              <Toggle
                labelPosition="right"
                checked={isActive}
                setChecked={() => {
                  if (!userData?.data.is_active) {
                    setIsActive((prev) => {
                      return !prev;
                    });
                  }
                }}
              >
                <p className="text-primary-black font-semibold">
                  {"Set as active address"}
                </p>
              </Toggle>
            </div>
          </div>
        </div>
        <div className="flex gap-4 justify-center">
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
                postalCode &&
                isDataChanged()
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
      {showNoActiveAddressErr &&
        createPortal(
          <ErrorBox onClose={() => setShowNoActiveAddressErr(false)}>
            <p className="font-bold text-primary-black text-4xl">Oops...</p>
            <p className="text-primary-black text-xl">
              An active address can't be deleted or deactivated before setting
              another address to be active
            </p>
          </ErrorBox>,
          document.body
        )}
      {showUpdateConfirm &&
        createPortal(
          <ConfirmBox
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
              if (!userData?.data.is_active) {
                handleDelete();
              } else {
                setShowNoActiveAddressErr(true);
              }
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
          isLoading || isDeleteLoading ? (
            <LoaderBox />
          ) : error || errDelete ? (
            <ErrorBox onClose={() => setShowResult(false)}>
              <p className="font-bold text-primary-black text-4xl">Oops...</p>
              <p className="text-primary-black text-xl">
                {(error && error[0].field === "server") ||
                (errDelete && errDelete[0].field === "server")
                  ? "Something is wrong, please try again"
                  : (error && CapitalizeFirstLetter(error[0].detail)) ||
                    (errDelete && CapitalizeFirstLetter(errDelete[0].detail))}
              </p>
            </ErrorBox>
          ) : (
            <SuccessBox onClose={() => navigate(PATH_BACK)}>
              Address {data ? "updated" : "deleted"} successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
