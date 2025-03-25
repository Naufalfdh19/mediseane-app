import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import { MapPin } from "lucide-react";

import Button from "../../../../components/Button";
import Input from "../../../../components/Input";
import ConfirmBox from "../../../../components/ConfirmBox";
import LoaderBox from "../../../../components/LoaderBox";
import ErrorBox from "../../../../components/ErrorBox";
import SuccessBox from "../../../../components/SuccessBox";
import {
  GetLocationByID,
  ParseIDAndNameFromLocation,
} from "../../../../utils/address";
import MapElem, { MapCoords } from "../../../../components/MapElem";
import ErrorCard from "../../../../components/ErrorCard";
import Dropdown from "../../../../components/Dropdown";
import TextArea from "../../../../components/TextArea";
import TitleWithBackButton from "../../../../components/TitleWithBackButton";
import DropdownHeader, {
  searchByType,
} from "../../../../components/DropdownHeader";

import {
  ADMIN_ADD_PHARMACIES_TITLE,
  API_ADDRESS_CHECK,
  API_ADDRESS_CITY,
  API_ADDRESS_DISTRICT,
  API_ADDRESS_PROVINCE,
  API_ADDRESS_SUBDISTRICT,
  API_ADMIN_PARTNERS,
  API_ADMIN_PHARMACY,
  API_METHOD_GET,
  API_METHOD_POST,
  PATH_ADMIN_PHARMACY,
} from "../../../../const/const";
import useAxios from "../../../../hooks/useAxios";
import {
  City,
  CurrLocation,
  District,
  Partner,
  PartnersData,
  Province,
  Response,
  SubDistrict,
} from "../../../../types/response";
import {
  GetPartnerByID,
  ParseIDAndNameFromPartner,
} from "../../../../utils/partner";
import { CapitalizeFirstLetter } from "../../../../utils/formatter";
import useTitle from "../../../../hooks/useTitle";

const searchByList: searchByType[] = [
  {
    id: "1",
    name: "name",
  },
  {
    id: "2",
    name: "year_founded",
  },
];

export default function CreatePharmacy() {
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState("");

  const [name, setName] = useState<string>();
  const [isNameValid, setIsNameValid] = useState<boolean>(true);
  const [partner, setPartner] = useState<Partner>();
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

  const navigate = useNavigate();

  const {
    data: partners,
    error: errPartner,
    isLoading: isLoadingPartner,
  } = useAxios<Response<PartnersData>>(
    API_ADMIN_PARTNERS +
      "?" +
      (searchByList[0] && filter ? `${searchByList[0].name}=${filter}` : ""),
    API_METHOD_GET
  );

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
    API_ADMIN_PHARMACY,
    API_METHOD_POST
  );

  useTitle(ADMIN_ADD_PHARMACIES_TITLE);

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

  function handleSubmit() {
    fetchData({
      name: name,
      partner_id: partner?.id,
      address: addrDetails,
      province_id: province?.id,
      province: province?.name,
      city_id: city?.id,
      city: city?.name,
      district_id: district?.id,
      district: district?.name,
      sub_district_id: subDistrict?.id,
      sub_district: subDistrict?.name,
      postal_code: postalCode ? parseInt(postalCode) : "",
      longitude: center?.lng.toString(),
      latitude: center?.lat.toString(),
    });
    setShowResult(true);
    setShowSaveModal(false);
  }

  function handleNameChange(e: React.ChangeEvent<HTMLInputElement>) {
    setName(e.target.value);
    setIsNameValid(e.target.value.length > 2);
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full">
          <TitleWithBackButton>Add Pharmacy</TitleWithBackButton>
        </div>
        <div className="w-full flex flex-col">
          {(errCurrLoc ||
            errPartner ||
            errProv ||
            errCity ||
            errDis ||
            errSubDis) && (
            <ErrorCard
              errors={
                errCurrLoc
                  ? errCurrLoc
                  : errPartner
                  ? errPartner
                  : errProv
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
        <div className="w-full">
          <div className="w-full flex flex-col gap-8">
            <div className="relative w-full flex flex-col gap-1.5">
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
                valid={isNameValid}
                onChange={handleNameChange}
              />
              {!isNameValid && (
                <p className="absolute bottom-[-26px] text-sm text-primary-red pl-1 pt-1">
                  Name should be at least 3 characters
                </p>
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="partner"
                className="text-primary-black font-semibold"
              >
                Partner <span className="text-primary-red">*</span>
              </label>
              <div className="w-full flex flex-col md:flex-row gap-6">
                <div className="md:w-[50%]">
                  <DropdownHeader
                    search={search}
                    setSearch={setSearch}
                    setFilter={setFilter}
                  />
                </div>
                <div className="md:w-[50%]">
                  <Dropdown
                    name="partner"
                    placeholder="Select partner"
                    data={ParseIDAndNameFromPartner(partners)}
                    selectedId={partner?.id.toString()}
                    onSelect={(id) => {
                      setPartner(GetPartnerByID(id, partners));
                    }}
                    disabled={errPartner !== undefined || isLoadingPartner}
                  />
                </div>
              </div>
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
        <div className="w-full flex justify-end">
          <div className="w-[367px] flex justify-between gap-4">
            <Button
              disabled={
                !(
                  name &&
                  partner &&
                  province &&
                  city &&
                  district &&
                  subDistrict &&
                  addrDetails &&
                  postalCode
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
                navigate(PATH_ADMIN_PHARMACY);
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
            <SuccessBox onClose={() => navigate(PATH_ADMIN_PHARMACY)}>
              Pharmacy added successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
