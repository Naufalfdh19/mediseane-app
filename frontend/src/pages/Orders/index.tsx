import { useEffect, useState } from "react";

import TitleWithBackButton from "../../components/TitleWithBackButton";
import OrderCard from "../../components/OrderCard";

import {
  API_USER_ORDERS,
  ORDER_STATUS_CANCELED,
  ORDER_STATUS_COMPLETED,
  ORDER_STATUS_DELIVERED,
  ORDER_STATUS_PENDING,
  ORDER_STATUS_PROCESSING,
  ORDER_STATUS_SHIPPED,
  ORDERS_TITLE,
} from "../../const/const";
import { Response, UserOrdersListItem } from "../../types/response";
import useAxios from "../../hooks/useAxios";
import ErrorCard from "../../components/ErrorCard";
import PageLoader from "../../components/PageLoader";
import { CapitalizeFirstLetter } from "../../utils/formatter";
import useTitle from "../../hooks/useTitle";
import ServerError from "../ServerError";

const types = [
  ORDER_STATUS_PENDING,
  ORDER_STATUS_PROCESSING,
  ORDER_STATUS_SHIPPED,
  ORDER_STATUS_DELIVERED,
  ORDER_STATUS_CANCELED,
];

export default function Orders() {
  const [page, setPage] = useState(0);
  const [orders, setOrders] = useState<UserOrdersListItem[]>([]);

  useTitle(ORDERS_TITLE);

  const { data, isLoading, error, fetchData } =
    useAxios<Response<UserOrdersListItem[]>>(API_USER_ORDERS);

  useEffect(() => {
    if (data && data.data) {
      setOrders(data.data);
    }
  }, [data]);

  function handlePage(thisPage: number) {
    if (page === thisPage) {
      return " bg-primary-green text-primary-white";
    }
    return " text-primary-green";
  }

  if (error && error[0].field === "server") {
    return <ServerError />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-7 md:my-16 w-[90%] max-w-[1108px] flex flex-col gap-7 md:gap-10 ">
        <TitleWithBackButton>Orders</TitleWithBackButton>
        <div className="w-full flex md:justify-center overflow-x-auto">
          <div className="rounded-full bg-secondary-green bg-opacity-[18%] flex p-1 justify-center items-center gap-2.5">
            {types.map((item, idx) => {
              return (
                <div
                  key={item}
                  className={
                    "flex justify-center items-center p-2.5 rounded-full cursor-pointer transition-all" +
                    handlePage(idx)
                  }
                  onClick={() => {
                    setPage(idx);
                  }}
                >
                  <p className="text-center">
                    {CapitalizeFirstLetter(
                      item.toLowerCase() === ORDER_STATUS_DELIVERED
                        ? ORDER_STATUS_COMPLETED
                        : item
                    )}
                  </p>
                </div>
              );
            })}
          </div>
        </div>
        {error ? (
          <ErrorCard errors={error} />
        ) : isLoading ? (
          <PageLoader />
        ) : (
          <div className="w-full flex flex-col gap-5">
            {orders
              .filter(
                (item) =>
                  item.order_details[0].status.toLowerCase() === types[page]
              )
              .map((item) => {
                return (
                  <OrderCard
                    key={item.id}
                    order={item}
                    refetchData={fetchData}
                  />
                );
              })}
            <div className="w-full flex justify-center mt-8">
              <p className="text-primary-gray text-center">
                This is the end of order history
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
