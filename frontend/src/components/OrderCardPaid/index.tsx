import { useState } from "react";
import { ChevronDown } from "lucide-react";
import { createPortal } from "react-dom";

import Button from "../Button";
import ConfirmBox from "../ConfirmBox";
import LoaderBox from "../LoaderBox";
import ErrorBox from "../ErrorBox";
import SuccessBox from "../SuccessBox";

import { Response, UserOrderDetails } from "../../types/response";
import { parseTimestamp } from "../../utils/timestamp";
import { statusColor } from "../../utils/orders";
import {
  API_METHOD_PATCH,
  API_USER_ORDERS,
  API_USER_ORDERS_SUFF_COMPLETION,
  ORDER_STATUS_COMPLETED,
  ORDER_STATUS_DELIVERED,
  ORDER_STATUS_SHIPPED,
} from "../../const/const";
import useAxios from "../../hooks/useAxios";
import { CapitalizeFirstLetter } from "../../utils/formatter";

interface OrderCardPropsPaid {
  order_date: string;
  order: UserOrderDetails;
  refetchData?: () => void;
}

export default function OrderCardPaid({
  order_date,
  order,
  refetchData = () => {},
}: OrderCardPropsPaid) {
  const [details, setDetails] = useState(false);
  const [showConfirmModal, setShowConfirmModal] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const {
    error: errConfirm,
    isLoading: isLoadingConfirm,
    fetchData: fetchConfirm,
  } = useAxios<Response<undefined>>(
    API_USER_ORDERS + "/" + order.details_id + API_USER_ORDERS_SUFF_COMPLETION,
    API_METHOD_PATCH
  );

  function handleConfirm() {
    fetchConfirm();
    setShowResult(true);
    setShowConfirmModal(false);
  }

  return (
    <div className="w-full bg-white shadow-md p-6 md:px-14 md:py-12 border border-primary-gray/20 rounded-lg flex flex-col gap-4 md:gap-10">
      <div className="flex flex-col gap-4 md:flex-row md:justify-between">
        <div className="flex flex-col gap-2.5">
          <p className="text-primary-black text-lg font-semibold">Pharmacy</p>
          <p className="text-primary-green font-bold">{order.pharmacy_name}</p>
        </div>
        <div className="flex flex-col gap-2.5">
          <p className="text-primary-black text-lg font-semibold">Order Date</p>
          <p className="text-primary-green font-bold">
            {parseTimestamp(order_date)}
          </p>
        </div>
        <div className="flex flex-col gap-2.5">
          <p className="text-primary-black text-lg font-semibold">
            Order Status
          </p>
          <p className={"font-bold" + " " + statusColor(order.status)}>
            {CapitalizeFirstLetter(
              order.status.toLowerCase() === ORDER_STATUS_DELIVERED
                ? ORDER_STATUS_COMPLETED
                : order.status
            )}
          </p>
        </div>
      </div>
      <div className="flex flex-col-reverse gap-4 md:flex-row md:justify-between md:items-center">
        <div
          className="flex gap-2.5 items-center text-primary-black cursor-pointer"
          onClick={() => {
            setDetails((prev) => !prev);
          }}
        >
          <ChevronDown
            className={`transform duration-500 ease-in-out ${
              details ? "rotate-180" : ""
            }`}
          />
          <p className="font-semibold text-lg">Details</p>
        </div>
        {order.status.toLowerCase() === ORDER_STATUS_SHIPPED && (
          <div className="w-full md:w-fit flex gap-2.5 justify-center flex-wrap items-center">
            <p className="text-primary-black font-bold">Confirm Reception:</p>
            <div>
              <Button
                size="sm"
                onClick={() => {
                  setShowConfirmModal(true);
                }}
              >
                <div className="px-4">
                  <p>Yes</p>
                </div>
              </Button>
            </div>
          </div>
        )}
      </div>
      {details && (
        <div className="flex flex-col gap-2.5">
          {order.order_products.map((product) => {
            return (
              <div
                key={product.order_product_id}
                className="flex flex-col md:flex-row md:flex-wrap md:justify-between md:items-center gap-6 py-2.5 border-t border-t-primary-gray/50"
              >
                <div className="flex justify-center gap-2.5 md:gap-10 items-center">
                  <div className=" relative w-[108px] h-[88px] md:w-[288px] md:h-[244px] flex justify-center items-center">
                    <img
                      className="object-contain"
                      src={product.image}
                      alt={product.name}
                    />
                  </div>
                  <div className="grow shrink flex flex-col gap-1 truncate md:max-w-[288px]">
                    <p className="text-primary-black text-lg truncate">
                      {product.name}
                    </p>
                    <p className="text-primary-gray text-base truncate">
                      {product.manufacturer}
                    </p>
                  </div>
                </div>
                <div className="flex flex-col gap-4 md:gap-6">
                  <div className="flex md:flex-col gap-2.5 items-center md:items-start truncate">
                    <p className="text-primary-black text-lg font-semibold">
                      Quantity
                    </p>
                    <p className="text-primary-green font-bold truncate">
                      {product.quantity}
                    </p>
                  </div>
                  <div className="flex md:flex-col gap-2.5 items-center md:items-start truncate">
                    <p className="text-primary-black text-lg font-semibold">
                      Subtotal
                    </p>
                    <p className="text-primary-green font-bold truncate">
                      {product.subtotal ? "Rp" + product.subtotal : ""}
                    </p>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
      {showConfirmModal &&
        createPortal(
          <ConfirmBox
            type="confirm"
            onYes={handleConfirm}
            onCancel={() => setShowConfirmModal(false)}
          />,
          document.body
        )}
      {showResult &&
        createPortal(
          isLoadingConfirm ? (
            <LoaderBox />
          ) : errConfirm ? (
            <ErrorBox onClose={() => setShowResult(false)}>
              <p className="font-bold text-primary-black text-4xl">Oops...</p>
              <p className="text-primary-black text-xl">
                {errConfirm && errConfirm[0].field === "server"
                  ? "Something is wrong, please try again"
                  : CapitalizeFirstLetter(errConfirm[0].detail)}
              </p>
            </ErrorBox>
          ) : (
            <SuccessBox
              onClose={() => {
                setShowResult(false);
                refetchData();
              }}
            >
              Order confirmed successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
