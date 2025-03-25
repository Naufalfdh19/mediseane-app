import { useEffect, useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { MapPinPlus } from "lucide-react";

import GroupedCheckoutCard from "./GroupedCheckoutCard";
import AddressCard from "../../components/AddressCard";
import TitleWithBackButton from "../../components/TitleWithBackButton";
import PageLoader from "../../components/PageLoader";
import Button from "../../components/Button";
import ErrorCard from "../../components/ErrorCard";
import ServerError from "../ServerError";

import useAxios from "../../hooks/useAxios";
import { useCheckoutState, useResetCheckout } from "../../store/checkout";
import {
  CheckoutDetails,
  DeliveryData,
  Response,
  UserAddress,
} from "../../types/response";
import {
  API_ADDRESS_USER,
  API_CHECKOUT,
  API_CHECKOUT_CONFIRM,
  API_METHOD_GET,
  API_METHOD_POST,
  CHECKOUT_TITLE,
  PATH_ADDRESS,
  PATH_USER_CART,
  PATH_USER_ORDERS,
} from "../../const/const";
import { useCart } from "../../hooks/useCart";
import useTitle from "../../hooks/useTitle";

export default function Checkout() {
  const [deliveries, setDeliveries] = useState<DeliveryData[]>();
  const [deliveryCosts, setDeliveryCosts] = useState<Map<number, number>>(
    new Map()
  );

  useTitle(CHECKOUT_TITLE);

  const { refetchCart } = useCart();

  const checkoutState = useCheckoutState();

  const resetCheckout = useResetCheckout();

  const navigate = useNavigate();

  const {
    data: address,
    isLoading: isAddressLoading,
    error: addressError,
  } = useAxios<Response<UserAddress[]>>(
    API_ADDRESS_USER + "?active=true",
    API_METHOD_GET
  );

  const { data, error, isLoading, fetchData } = useAxios<
    Response<CheckoutDetails>
  >(API_CHECKOUT, API_METHOD_POST);

  const {
    error: checkoutError,
    isLoading: isCheckoutLoading,
    fetchData: postCheckout,
  } = useAxios<Response<undefined>>(API_CHECKOUT_CONFIRM, API_METHOD_POST);

  useEffect(() => {
    fetchData({ ids: checkoutState });
  }, []);

  useEffect(() => {
    if (data) {
      let deliveryDatas: DeliveryData[] = [];
      data.data.data.forEach((item) => {
        deliveryDatas.push({
          pharmacy_id: item.pharmacy_id,
          delivery_id: 0,
        });
      });
      setDeliveries(deliveryDatas);
    }
  }, [data]);

  function isCheckoutValid() {
    if (deliveries) {
      return (
        deliveries.filter((delivery) => delivery.delivery_id === 0).length > 0
      );
    }
  }

  function calculateTotalShipping() {
    if (deliveryCosts) {
      let total: number = 0;
      Array.from(deliveryCosts.entries()).map(([_, value]) => {
        total += +value;
      });
      return total;
    }
  }

  function calculateTotalProducts() {
    let total: number = 0;
    data?.data.data.forEach((item) =>
      item.items.forEach((item) => (total += +item.subtotal))
    );
    return total;
  }

  function handleCheckout() {
    postCheckout({
      id_cart: data?.data.id,
      delivery_data_list: deliveries,
    }).then((res) => {
      if (res && res.message) {
        refetchCart();
        resetCheckout();
        navigate(PATH_USER_ORDERS);
      }
    });
  }

  if (checkoutState.length < 1) {
    return <Navigate to={PATH_USER_CART} />;
  }

  if (isLoading || isAddressLoading) {
    return <PageLoader />;
  }

  if (
    (error && error[0].field === "server") ||
    (addressError && addressError[0].field === "server")
  ) {
    return <ServerError />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-7 md:my-16 w-[90%] max-w-[1108px] flex flex-col gap-7 md:gap-10">
        <TitleWithBackButton>Checkout</TitleWithBackButton>
        {error ? (
          <ErrorCard errors={error} noMargin={true} />
        ) : (
          addressError && <ErrorCard errors={addressError} noMargin={true} />
        )}
        <div className="flex flex-col gap-4">
          {address && address.data.length < 1 ? (
            <div
              className="w-full h-[170px] bg-primary-green bg-opacity-[18%] cursor-pointer shadow-md p-4 border border-primary-gray/20 rounded-lg flex gap-2 justify-center items-center"
              onClick={() => navigate(PATH_ADDRESS)}
            >
              <p className="text-primary-green">
                Click here to set{" "}
                <span className="font-bold">your address</span>
              </p>
              <MapPinPlus
                strokeWidth={1.4}
                className="text-primary-green flex-shrink-0"
              />
            </div>
          ) : (
            <div className="relative" onClick={() => navigate(PATH_ADDRESS)}>
              <AddressCard address={address!.data[0]} />
              <p className="absolute top-5 right-5 text-primary-green font-semibold cursor-pointer">
                Change
              </p>
            </div>
          )}
          <div className="flex flex-col md:flex-row gap-3">
            <div className="w-full md:w-[65%] flex flex-col gap-4">
              {data?.data &&
                data.data.data &&
                data.data.data.map((item, idx) => (
                  <GroupedCheckoutCard
                    key={item.pharmacy_id}
                    item={item}
                    isAddressSet={address!.data.length > 0}
                    shippingID={deliveries ? deliveries[idx].delivery_id : 0}
                    deliveries={deliveries}
                    setDeliveries={setDeliveries}
                    deliveryCosts={deliveryCosts}
                    setDeliveryCosts={setDeliveryCosts}
                  />
                ))}
            </div>

            <div className="md:w-[35%] h-fit bg-white shadow-md p-4 border border-primary-gray/20 rounded-lg flex flex-col gap-6">
              <div className="flex flex-col gap-1 overflow-auto">
                <p className="text-primary-black font-semibold">
                  Payment Details
                </p>
                <div className="flex justify-between gap-2">
                  <p className="text-primary-black">Subtotal products: </p>
                  <p className="text-primary-green">
                    Rp{calculateTotalProducts().toFixed(2)}
                  </p>
                </div>
                <div className="flex justify-between gap-2">
                  <p className="text-primary-black">Subtotal shipping cost: </p>
                  <p className="text-primary-green">
                    Rp{calculateTotalShipping()}
                  </p>
                </div>
                <div className="flex justify-between gap-2">
                  <p className="text-primary-black">Total: </p>
                  <p className="text-primary-green font-semibold">
                    Rp{calculateTotalProducts()! + calculateTotalShipping()!}
                  </p>
                </div>
              </div>

              <Button
                disabled={isCheckoutValid() || isCheckoutLoading}
                onClick={handleCheckout}
              >
                {isCheckoutLoading ? "Loading..." : "Confirm"}
              </Button>

              {checkoutError && (
                <ErrorCard errors={checkoutError} noMargin={true} />
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
