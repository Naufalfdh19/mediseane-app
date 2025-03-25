import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { createPortal } from "react-dom";

import Input from "../../../../components/Input";
import Button from "../../../../components/Button";
import ConfirmBox from "../../../../components/ConfirmBox";
import LoaderBox from "../../../../components/LoaderBox";
import ErrorBox from "../../../../components/ErrorBox";
import SuccessBox from "../../../../components/SuccessBox";
import { CapitalizeFirstLetter } from "../../../../utils/formatter";

import { Response } from "../../../../types/response";
import useAxios from "../../../../hooks/useAxios";
import {
  API_METHOD_GET,
  PATH_BACK,
  API_METHOD_PATCH,
  PHARMACIST_EDIT_PRODUCT_TITLE,
} from "../../../../const/const";

import { PharmacyProduct } from "../../../../types/response";
import TitleWithBackButton from "../../../../components/TitleWithBackButton";
import ErrorCard from "../../../../components/ErrorCard";
import NotFound from "../../../NotFound";
import Toggle from "../../../../components/Toggle";
import {
  API_PHARMACIST_PRODUCTS,
  PATH_PHARMACIST_PRODUCTS,
} from "../../../../const/const";
import useTitle from "../../../../hooks/useTitle";

export default function EditProduct() {
  const [stock, setStock] = useState<number>();
  const [isStockValid, setIsStockValid] = useState<boolean>(true);
  const [isActive, setIsActive] = useState<boolean>(false);
  const [name, setName] = useState<string>();
  const [genericName, setGenericName] = useState<string>();
  const [manufacturer, setManufacturer] = useState<string>();
  const [classification, setClassification] = useState<string>();
  const [form, setForm] = useState<string>();
  const [price, setPrice] = useState<string>();
  const [fullPharmacyProduct, setFullPharmacyProduct] =
    useState<PharmacyProduct>();

  const [showUpdateConfirm, setShowUpdateConfirm] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const navigate = useNavigate();

  const params = useParams();

  const {
    data: pharmacy_products,
    error: errPharmacyProduct,
    isLoading: isLoadingPharmacyProduct,
  } = useAxios<Response<PharmacyProduct>>(
    API_PHARMACIST_PRODUCTS + "/" + params.id,
    API_METHOD_GET
  );

  const { data, error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_PHARMACIST_PRODUCTS + "/" + params.id,
    API_METHOD_PATCH
  );

  useTitle(PHARMACIST_EDIT_PRODUCT_TITLE);

  useEffect(() => {
    if (pharmacy_products?.data) {
      setFullPharmacyProduct(pharmacy_products.data);
    }
  }, [pharmacy_products]);

  useEffect(() => {
    if (fullPharmacyProduct) {
      setName(fullPharmacyProduct.name);
      setGenericName(fullPharmacyProduct.generic_name);
      setManufacturer(fullPharmacyProduct.manufacturer);
      setClassification(fullPharmacyProduct.product_classification);
      setForm(fullPharmacyProduct.product_form);
      setStock(fullPharmacyProduct.stock);
      setPrice(fullPharmacyProduct.price);
      setIsActive(fullPharmacyProduct.is_active);
    }
  }, [fullPharmacyProduct]);

  function handleStockChange(e: React.ChangeEvent<HTMLInputElement>) {
    const stock = Number(e.target.value);
    setStock(stock);
    setIsStockValid(stock > 0);
  }

  function handleSubmit() {
    fetchData({
      id: params.id ? parseInt(params.id) : undefined,
      stock: stock,
      is_active: isActive,
    });
    setShowResult(true);
  }

  if (errPharmacyProduct && errPharmacyProduct[0].field === "not found") {
    return <NotFound />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-7 md:my-16 w-[90%] max-w-[1108px] flex flex-col gap-7 md:gap-10">
        <TitleWithBackButton>Edit Product</TitleWithBackButton>
        <div className="w-full">
          {errPharmacyProduct && (
            <ErrorCard errors={errPharmacyProduct ? errPharmacyProduct : []} />
          )}
          <div className="w-full">
            <div className="w-full flex flex-col gap-8">
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Product name
                </label>
                {isLoadingPharmacyProduct ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="name"
                    placeholder=""
                    type="text"
                    value={name ? name : ""}
                    readOnly
                  />
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="stock"
                  className="text-primary-black font-semibold"
                >
                  Stock
                </label>
                {isLoadingPharmacyProduct ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="stock"
                    placeholder=""
                    type="number"
                    value={stock?.toString() ? stock?.toString() : ""}
                    valid={isStockValid}
                    onChange={handleStockChange}
                  />
                )}
                {!isStockValid && (
                  <p className="text-sm text-primary-red pl-1 pt-1">
                    Stock should be more than 0
                  </p>
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
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
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Generic name
                </label>
                {isLoadingPharmacyProduct ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="generic"
                    placeholder=""
                    type="text"
                    value={genericName ? genericName : ""}
                    readOnly
                  />
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Manufacturer
                </label>
                {isLoadingPharmacyProduct ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="manufacturer"
                    placeholder=""
                    type="text"
                    value={manufacturer ? manufacturer : ""}
                    readOnly
                  />
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Product classification
                </label>
                {isLoadingPharmacyProduct ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="classification"
                    placeholder=""
                    type="text"
                    value={classification ? classification : ""}
                    readOnly
                  />
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Product form
                </label>
                {isLoadingPharmacyProduct ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="form"
                    placeholder=""
                    type="text"
                    value={form ? form : ""}
                    readOnly
                  />
                )}
              </div>
              <div className="w-full flex flex-col gap-1.5">
                <label
                  htmlFor="name"
                  className="text-primary-black font-semibold"
                >
                  Price
                </label>
                {isLoadingPharmacyProduct ? (
                  <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                ) : (
                  <Input
                    name="price"
                    placeholder=""
                    type="text"
                    value={price ? "Rp" + price : ""}
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
              disabled={!stock || isLoading}
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
            <SuccessBox onClose={() => navigate(PATH_PHARMACIST_PRODUCTS)}>
              Pharmacy product {data ? "updated" : "deleted"} successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
