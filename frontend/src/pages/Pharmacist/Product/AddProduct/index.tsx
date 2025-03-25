import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";

import Button from "../../../../components/Button";
import Input from "../../../../components/Input";
import ConfirmBox from "../../../../components/ConfirmBox";
import LoaderBox from "../../../../components/LoaderBox";
import ErrorBox from "../../../../components/ErrorBox";
import SuccessBox from "../../../../components/SuccessBox";
import Dropdown from "../../../../components/Dropdown";

import {
  API_MASTER_PRODUCTS,
  API_METHOD_GET,
  API_METHOD_POST,
  API_PHARMACIST_PRODUCTS,
  PATH_PHARMACIST_PRODUCTS,
  PHARMACIST_ADD_PRODUCT_TITLE,
} from "../../../../const/const";
import useAxios from "../../../../hooks/useAxios";
import {
  ProductListItem,
  ProductResponse,
  Response,
} from "../../../../types/response";
import {
  GetProductByID,
  ParseIDAndNameFromProduct,
} from "../../../../utils/product";
import TitleWithBackButton from "../../../../components/TitleWithBackButton";
import DropdownHeader, {
  searchByType,
} from "../../../../components/DropdownHeader";
import Toggle from "../../../../components/Toggle";
import { CapitalizeFirstLetter } from "../../../../utils/formatter";
import useTitle from "../../../../hooks/useTitle";

const searchByList: searchByType[] = [
  {
    id: "1",
    name: "name",
  },
];

export default function AddProduct() {
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState("");

  const [product, setProduct] = useState<ProductListItem>();
  const [stock, setStock] = useState<number>();
  const [isStockValid, setIsStockValid] = useState<boolean>(true);
  const [price, setPrice] = useState<string>();
  const [isPriceValid, setIsPriceValid] = useState<boolean>(true);
  const [isActive, setIsActive] = useState<boolean>(false);
  const [showSaveModal, setShowSaveModal] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const navigate = useNavigate();

  const {
    data: products,
    error: errProduct,
    isLoading: isLoadingProduct,
  } = useAxios<Response<ProductResponse>>(
    API_MASTER_PRODUCTS +
      "?" +
      (searchByList[0] && filter ? `${searchByList[0].name}=${filter}` : ""),
    API_METHOD_GET
  );

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_PHARMACIST_PRODUCTS,
    API_METHOD_POST
  );

  useTitle(PHARMACIST_ADD_PRODUCT_TITLE)

  function handlePriceChange(e: React.ChangeEvent<HTMLInputElement>) {
    const price = e.target.value;
    setPrice(price);
    setIsPriceValid(Number(price) > 0);
  }

  function handleStockChange(e: React.ChangeEvent<HTMLInputElement>) {
    const stock = Number(e.target.value);
    setStock(stock);
    setIsStockValid(stock > 0);
  }

  function handleSubmit() {
    fetchData({
      product_id: product?.id,
      stock: stock,
      price: price,
      is_active: isActive,
    });
    setShowResult(true);
    setShowSaveModal(false);
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full">
          <TitleWithBackButton>Add Product</TitleWithBackButton>
        </div>
        <div className="w-full">
          <div className="w-full flex flex-col gap-8">
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="product"
                className="text-primary-black font-semibold"
              >
                Product <span className="text-primary-red">*</span>
              </label>
              <div className="w-full flex flex-col md:flex-row gap-6">
                <div className="md:w-[50%]">
                  {isLoadingProduct ? (
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
                  {isLoadingProduct ? (
                    <div className="w-full h-12 gap-2 rounded-lg bg-secondary-gray animate-pulse"></div>
                  ) : (
                    <Dropdown
                      name="product"
                      placeholder="Select product"
                      data={ParseIDAndNameFromProduct(products)}
                      selectedId={product?.id.toString()}
                      onSelect={(id) => {
                        setProduct(GetProductByID(id, products));
                      }}
                      disabled={errProduct !== undefined || isLoadingProduct}
                    />
                  )}
                </div>
              </div>
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="name"
                className="text-primary-black font-semibold"
              >
                Stock <span className="text-primary-red">*</span>
              </label>
              <Input
                name="stock"
                placeholder="1"
                type="number"
                valid={isStockValid}
                onChange={handleStockChange}
              />
              {!isStockValid && (
                <p className="text-sm text-primary-red pl-1 pt-1">
                  Stock should be more than 0
                </p>
              )}
            </div>
            <div className="w-full flex flex-col gap-1.5">
              <label
                htmlFor="name"
                className="text-primary-black font-semibold"
              >
                Price <span className="text-primary-red">*</span>
              </label>
              <Input
                name="price"
                placeholder="1"
                type="number"
                valid={isPriceValid}
                onChange={handlePriceChange}
              />
              {!isPriceValid && (
                <p className="text-sm text-primary-red pl-1 pt-1">
                  Price should be more than 0
                </p>
              )}
            </div>
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
        </div>
        <div className="w-full flex justify-end">
          <div className="w-[367px] flex justify-between gap-4">
            <Button
              disabled={!(product && stock && price) || isLoading}
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
                navigate(PATH_PHARMACIST_PRODUCTS);
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
                  : error && CapitalizeFirstLetter(error[0].detail)}
              </p>
            </ErrorBox>
          ) : (
            <SuccessBox onClose={() => navigate(PATH_PHARMACIST_PRODUCTS)}>
              Product added successfully to the pharmacy!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
