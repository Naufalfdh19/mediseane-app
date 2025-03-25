import { useNavigate } from "react-router-dom";

import GroupedCartCard from "./GroupedCartCard";
import Button from "../../components/Button";
import TitleWithBackButton from "../../components/TitleWithBackButton";
import ErrorCard from "../../components/ErrorCard";
import PageLoader from "../../components/PageLoader";
import ServerError from "../ServerError";

import useAxios from "../../hooks/useAxios";
import { useCheckoutState } from "../../store/checkout";
import { GroupedCartItem, Response } from "../../types/response";
import {
  API_CART,
  API_METHOD_GET,
  CARTS_TITLE,
  PATH_CHECKOUT,
  PATH_PRODUCTS,
} from "../../const/const";
import { ShoppingCart } from "lucide-react";
import useTitle from "../../hooks/useTitle";

export default function Cart() {
  useTitle(CARTS_TITLE);

  const checkoutState = useCheckoutState();
  const { data, error, isLoading, fetchData } = useAxios<
    Response<GroupedCartItem[]>
  >(API_CART, API_METHOD_GET);

  const navigate = useNavigate();

  function countTotal() {
    let total = 0;
    data?.data.forEach((pharmacy) => {
      pharmacy.items.filter((item) => {
        if (checkoutState.includes(item.cart_item_id)) {
          total += +item.subtotal;
        }
      });
    });
    return total;
  }

  if (isLoading && !data) {
    return <PageLoader />;
  }

  if (error && error[0].field === "server") {
    return <ServerError />;
  }

  return (
    <div className="grow flex flex-col bg-primary-white">
      <div className="my-7 md:my-16 mx-auto w-[90%] max-w-[627px] flex flex-col gap-7 md:gap-10">
        <TitleWithBackButton>Cart</TitleWithBackButton>
        <div className="w-full flex flex-col gap-5">
          {error && <ErrorCard errors={error} />}
          {data && data.data.length < 1 ? (
            <div className="h-[60%] flex flex-col gap-3 items-center justify-center">
              <ShoppingCart
                className="h-32 w-32 text-primary-green"
                strokeWidth={0.8}
              />
              <p className="text-lg text-primary-black">
                It looks like your cart is empty
              </p>
              <div className="w-full max-w-[300px]">
                <Button onClick={() => navigate(PATH_PRODUCTS)}>
                  Shop now
                </Button>
              </div>
            </div>
          ) : (
            data?.data &&
            data.data.map((items) => {
              return (
                <GroupedCartCard
                  key={items.pharmacy_id}
                  items={items}
                  refetch={fetchData}
                />
              );
            })
          )}
        </div>
        <div className="sticky bottom-7 bg-white cursor-pointer shadow-md p-4 border border-primary-gray/20 rounded-lg flex items-center gap-3">
          <div className="w-[50%]">
            Total price:{" "}
            <span className="text-primary-green font-semibold">
              Rp{countTotal().toFixed(2)}
            </span>
          </div>
          <div className="w-[50%]">
            <Button
              disabled={checkoutState.length < 1}
              onClick={() => navigate(PATH_CHECKOUT)}
            >
              Checkout {checkoutState.length > 0 && `(${checkoutState.length})`}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
