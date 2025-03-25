import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { createPortal } from "react-dom";
import { ShoppingCart } from "lucide-react";

import SearchHeader from "../../components/SearchHeader";
import PageLoader from "../../components/PageLoader";
import ErrorBox from "../../components/ErrorBox";
import Button from "../../components/Button";
import Tag from "../../components/Tag";
import AddToCartButton from "../../components/AddToCartButton";
import Modal from "../../components/Modal/Modal";
import ServerError from "../ServerError";

import { Response, UserProductDetails } from "../../types/response";
import { CapitalizeFirstLetter } from "../../utils/formatter";
import { useUserState } from "../../store/user";
import useAxios from "../../hooks/useAxios";
import {
  API_METHOD_GET,
  API_USER_PRODUCT_DETAILS,
  PATH_LOGIN,
  PATH_PROFILE,
  PRODUCTS_TITLE,
} from "../../const/const";
import useTitle from "../../hooks/useTitle";

export default function ProductDetails() {
  const [isNotAbleToShop, setIsNotAbleToShop] = useState<boolean>(false);

  const params = useParams();
  const navigate = useNavigate();
  const userState = useUserState();

  const id = params.id;

  const { data, isLoading, error } = useAxios<Response<UserProductDetails>>(
    API_USER_PRODUCT_DETAILS + id,
    API_METHOD_GET
  );

  useTitle(data ? data.data.name : PRODUCTS_TITLE);

  function handleNotLogin() {
    navigate(PATH_LOGIN);
  }

  function handleNotVerified() {
    navigate(PATH_PROFILE);
  }

  if (isLoading) {
    return <PageLoader />;
  }

  if (error && error[0].field === "server") {
    return <ServerError />;
  }

  return (
    <>
      <SearchHeader />
      <div className="grow flex bg-primary-white justify-center">
        <div className="my-16 w-[80%] md:w-[90%] max-w-[1259px] flex flex-col md:flex-row justify-start md:justify-center gap-10">
          <div className="bg-pure-white h-[245px] md:h-[523px] w-full md:w-[59%] flex justify-center items-center">
            <img
              className="object-contain w-full h-full"
              src={data?.data.image.split("'")[1]}
              alt={data?.data.name}
            />
          </div>
          <div className="w-full md:w-[39%] flex flex-col gap-8 md:gap-10">
            <div className="w-full flex flex-col gap-5">
              <div className="w-full flex flex-col gap-3.5">
                <h1 className="text-2xl text-primary-black">
                  {data?.data.name}
                </h1>
                <h1 className="text-primary-green text-lg md:text-[32px] font-semibold">
                  {data?.data.price ? "Rp" + data?.data.price : ""}
                </h1>
              </div>
              <div className="flex flex-wrap gap-2.5">
                {data?.data &&
                  data.data.product_categories &&
                  data.data.product_categories.map((item) => {
                    return <Tag key={item}>{item}</Tag>;
                  })}
              </div>
            </div>
            <div className="flex flex-col gap-2.5">
              <p className="text-primary-gray text-base font-bold">
                Generic Name:{" "}
                <span className="font-normal">{data?.data.generic_name}</span>
              </p>
              <p className="text-primary-gray text-base font-bold">
                Manufacturer:{" "}
                <span className="font-normal">{data?.data.manufacture}</span>
              </p>
              <p className="text-primary-gray text-base font-bold">
                Description:{" "}
                <span className="font-normal">{data?.data.description}</span>
              </p>
              <p className="text-primary-gray text-base font-bold">
                Unit in Pack:{" "}
                <span className="font-normal">{data?.data.unit_in_pack}</span>
              </p>
            </div>
            <div className="flex flex-col gap-2.5">
              <p className="text-primary-gray text-base font-bold">
                Pharmacy Name:{" "}
                <span className="font-normal">
                  {data?.data.pharmacies_name}
                </span>
              </p>
              <p className="text-primary-gray text-base font-bold">
                Pharmacy Address:{" "}
                <span className="font-normal">{data?.data.address}</span>
              </p>
            </div>
            <div className="flex flex-col gap-3.5">
              <p className="text-primary-gray text-base font-bold">
                Available Stock:{" "}
                <span className="font-normal">{data?.data.stock}</span>
              </p>
            </div>
            <div className="w-full sticky bottom-7 p-4 bg-primary-white rounded-lg flex items-center gap-3">
              {userState.id !== 0 && userState.is_verified ? (
                <AddToCartButton
                  product_id={data ? data.data.pharmacy_product_id : 0}
                />
              ) : (
                <Button onClick={() => setIsNotAbleToShop(true)}>
                  <div className="flex justify-center items-center gap-1.5">
                    <ShoppingCart />
                    <h1>Add to Cart</h1>
                  </div>
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>
      {error &&
        createPortal(
          <ErrorBox
            onClose={() => {
              navigate(-1);
            }}
          >
            <p className="font-bold text-primary-black text-4xl">Oops...</p>
            <p className="text-primary-black text-xl">
              {error && error[0].field === "server"
                ? "Something is wrong, please try again"
                : CapitalizeFirstLetter(error[0].detail)}
            </p>
          </ErrorBox>,
          document.body
        )}
      {isNotAbleToShop &&
        createPortal(
          <Modal onClose={() => setIsNotAbleToShop(false)}>
            <div className="w-full py-8 px-2 flex flex-col items-center justify-center gap-6">
              <ShoppingCart className="w-[100px] h-[100px] md:w-[130px] md:h-[130px] text-primary-green" />
              <p className="px-2 font-semibold text-primary-black text-2xl md:text-3xl text-center">
                {userState.id === 0 ? "Login" : "Verify email"} before adding to
                cart
              </p>
              <div className="w-full flex flex-wrap gap-3 md:gap-7 justify-center items-center">
                <div className="w-[120px] md:w-[163px]">
                  <Button
                    submit={false}
                    size="md"
                    onClick={
                      userState.id === 0 ? handleNotLogin : handleNotVerified
                    }
                  >
                    {userState.id === 0 ? "Login now" : "Verify email"}
                  </Button>
                </div>
                <div className="w-[120px] md:w-[163px]">
                  <Button
                    submit={false}
                    type="ghost"
                    size="md"
                    onClick={() => setIsNotAbleToShop(false)}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            </div>
          </Modal>,
          document.body
        )}
    </>
  );
}
