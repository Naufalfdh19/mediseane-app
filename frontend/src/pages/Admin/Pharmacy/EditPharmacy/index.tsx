import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { createPortal } from "react-dom";

import MapElem, { MapCoords } from "../../../../components/MapElem";
import NotFound from "../../../NotFound";
import TitleWithBackButton from "../../../../components/TitleWithBackButton";
import ErrorCard from "../../../../components/ErrorCard";
import Input from "../../../../components/Input";
import Dropdown from "../../../../components/Dropdown";
import TextArea from "../../../../components/TextArea";
import Button from "../../../../components/Button";
import LoaderBox from "../../../../components/LoaderBox";
import ConfirmBox from "../../../../components/ConfirmBox";
import ErrorBox from "../../../../components/ErrorBox";
import { CapitalizeFirstLetter } from "../../../../utils/formatter";
import SuccessBox from "../../../../components/SuccessBox";
import DropdownHeader, {
  searchByType,
} from "../../../../components/DropdownHeader";
import ServerError from "../../../ServerError";

import {
  City,
  District,
  Partner,
  PartnersData,
  PharmacyItemWithID,
  Province,
  Response,
  SubDistrict,
} from "../../../../types/response";
import useAxios from "../../../../hooks/useAxios";
import {
  ADMIN_EDIT_PHARMACIES_TITLE,
  API_ADDRESS_CITY,
  API_ADDRESS_DISTRICT,
  API_ADDRESS_PROVINCE,
  API_ADDRESS_SUBDISTRICT,
  API_ADMIN_PARTNERS,
  API_ADMIN_PHARMACY,
  API_METHOD_GET,
  API_METHOD_PUT,
  PATH_ADMIN_PHARMACY,
} from "../../../../const/const";
import {
  GetPartnerByID,
  ParseIDAndNameFromPartner,
} from "../../../../utils/partner";
import {
  GetLocationByID,
  ParseIDAndNameFromLocation,
} from "../../../../utils/address";
import useTitle from "../../../../hooks/useTitle";
import Toggle from "../../../../components/Toggle";

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

export default function EditPharmacy() {
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState("");

  const [name, setName] = useState<string>("");
  const [isActive, setIsActive] = useState<boolean>(false);
  const [partner, setPartner] = useState<Partner>();
  const [province, setProvince] = useState<Province>();
  const [city, setCity] = useState<City>();
  const [district, setDistrict] = useState<District>();
  const [subDistrict, setSubDistrict] = useState<SubDistrict>();
  const [addrDetails, setAddrDetails] = useState<string>("");
  const [postalCode, setPostalCode] = useState<string>("");
  const [center, setCenter] = useState<MapCoords>();
  const [fullPharmacy, setFullPharmacy] = useState<PharmacyItemWithID>();
  const [showResult, setShowResult] = useState<boolean>(false);
  const [showUpdateConfirm, setShowUpdateConfirm] = useState<boolean>(false);
  const [showNoActiveAddressErr, setShowNoActiveAddressErr] =
    useState<boolean>(false);

  const navigate = useNavigate();

  const params = useParams();

  const {
    data: pharmacyData,
    error: errPharmacyData,
    isLoading: isLoadingPharmacyData,
  } = useAxios<Response<PharmacyItemWithID>>(
    API_ADMIN_PHARMACY + "/" + params.id,
    API_METHOD_GET
  );

  const {
    data: partners,
    error: errPartner,
    isLoading: isLoadingPartner,
    fetchData: fetchPartner,
  } = useAxios<Response<PartnersData>>(
    API_ADMIN_PARTNERS +
      "?" +
      (searchByList[0] && filter ? `${searchByList[0].name}=${filter}` : ""),
    API_METHOD_GET,
    false,
    true
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
    API_ADMIN_PHARMACY,
    API_METHOD_PUT
  );

  useTitle(ADMIN_EDIT_PHARMACIES_TITLE);

  useEffect(() => {
    if (pharmacyData?.data) {
      setFullPharmacy(pharmacyData.data);
    }
  }, [pharmacyData]);

  useEffect(() => {
    if (fullPharmacy) {
      setName(fullPharmacy.name);
      setIsActive(fullPharmacy.is_active);
      setAddrDetails(fullPharmacy.address);
      fetchProvince();
      setSearch(fullPharmacy.partner_name);
      setFilter(fullPharmacy.partner_name);
    }
  }, [fullPharmacy]);

  useEffect(() => {
    fetchPartner();
  }, [filter]);

  useEffect(() => {
    if (partners && !partner && fullPharmacy?.partner_id) {
      setPartner(GetPartnerByID(fullPharmacy.partner_id.toString(), partners));
    }
  }, [partners]);

  useEffect(() => {
    if (provinces && !province && fullPharmacy?.province_id) {
      setProvince(
        GetLocationByID(fullPharmacy.province_id.toString(), provinces)
      );
    }
  }, [provinces]);

  useEffect(() => {
    if (cities && !city && fullPharmacy?.city_id) {
      setCity(GetLocationByID(fullPharmacy.city_id.toString(), cities));
    }
  }, [cities]);

  useEffect(() => {
    if (districts && !district && fullPharmacy?.district_id) {
      setDistrict(
        GetLocationByID(fullPharmacy.district_id.toString(), districts)
      );
    }
  }, [districts]);

  useEffect(() => {
    if (subDistricts && !subDistrict && fullPharmacy?.sub_district_id) {
      setSubDistrict(
        GetLocationByID(fullPharmacy.sub_district_id.toString(), subDistricts)
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
    if (fullPharmacy) {
      setPostalCode(fullPharmacy.postal_code.toString());
      setCenter({
        lat: parseFloat(fullPharmacy.latitude),
        lng: parseFloat(fullPharmacy.longitude),
      });
      setFullPharmacy(undefined);
    }
  }, [subDistrict]);

  function handleSubmit() {
    fetchData({
      id: params.id ? parseInt(params.id) : undefined,
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
      is_active: isActive,
      longitude: center?.lng.toString(),
      latitude: center?.lat.toString(),
    });
    setShowResult(true);
    setShowUpdateConfirm(false);
  }

  if (errPharmacyData && errPharmacyData[0].field === "not found") {
    return <NotFound />;
  }

  if (errPharmacyData && errPharmacyData[0].field === "server") {
    return <ServerError />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-7 md:my-16 w-[90%] max-w-[1108px] flex flex-col gap-7 md:gap-10">
        <TitleWithBackButton>Edit Pharmacy</TitleWithBackButton>
        <div className="w-full">
          {(errPartner || errProv || errCity || errDis || errSubDis) && (
            <ErrorCard
              errors={
                errPartner
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
          <div className="w-full flex flex-col gap-8">
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="name"
                className="text-primary-black font-semibold"
              >
                Name
              </label>
              {isLoadingPartner ? (
                <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
              ) : (
                <Input
                  name="name"
                  placeholder="Name"
                  type="text"
                  value={name}
                  onChange={(e) => {
                    setName(e.target.value);
                  }}
                />
              )}
            </div>
            <div>
              <Toggle
                labelPosition="right"
                checked={isActive}
                setChecked={() => {
                  setIsActive((prev) => {
                    return !prev;
                  });
                }}
              >
                <p className="text-primary-black font-semibold">{"Active"}</p>
              </Toggle>
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="partner"
                className="text-primary-black font-semibold"
              >
                Partner
              </label>
              <div className="w-full flex flex-col md:flex-row gap-6">
                <div className="w-full flex flex-col md:flex-row gap-6">
                  <div className="md:w-[50%]">
                    {isLoadingPartner ? (
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
                    {isLoadingPartner ? (
                      <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                    ) : (
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
                    )}
                  </div>
                </div>
              </div>
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="province"
                className="text-primary-black font-semibold"
              >
                Province
              </label>
              {isLoadingProv ? (
                <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
              ) : (
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
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="city"
                className="text-primary-black font-semibold"
              >
                City
              </label>
              {isLoadingCity ? (
                <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
              ) : (
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
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="district"
                className="text-primary-black font-semibold"
              >
                District
              </label>
              {isLoadingDis ? (
                <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
              ) : (
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
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="sub-district"
                className="text-primary-black font-semibold"
              >
                Sub district
              </label>
              {isLoadingSubDis ? (
                <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
              ) : (
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
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="postcode"
                className="text-primary-black font-semibold"
              >
                Postal code
              </label>
              {isLoadingPharmacyData ? (
                <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
              ) : (
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
              )}
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
                setShowUpdateConfirm(true);
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
            <SuccessBox onClose={() => navigate(PATH_ADMIN_PHARMACY)}>
              Pharmacy {data ? "updated" : "deleted"} successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
