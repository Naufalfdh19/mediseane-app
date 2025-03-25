import { Camera, ChevronDown, Hospital, Image } from "lucide-react";
import { useRef, useState } from "react";
import { createPortal } from "react-dom";

import Button from "../Button";
import OrderCardPaid from "../OrderCardPaid";
import ConfirmBox from "../ConfirmBox";
import LoaderBox from "../LoaderBox";
import ErrorBox from "../ErrorBox";
import SuccessBox from "../SuccessBox";
import Modal from "../Modal/Modal";

import {
  API_METHOD_DELETE,
  API_METHOD_PATCH,
  API_USER_ORDER_DETAILS,
  API_USER_ORDER_DETAILS_SUFF_PAYMENT,
  API_USER_ORDERS,
  ORDER_STATUS_PENDING,
  ORDER_STATUS_WAITING,
} from "../../const/const";
import { Response, UserOrdersListItem } from "../../types/response";
import { statusColor } from "../../utils/orders";
import { parseTimestamp } from "../../utils/timestamp";
import useAxios from "../../hooks/useAxios";
import { CapitalizeFirstLetter } from "../../utils/formatter";
const types = ["jpg", "jpeg", "png", "pdf"];

interface OrderCardProps {
  order: UserOrdersListItem;
  refetchData?: () => void;
}

export default function OrderCard({
  order,
  refetchData = () => {},
}: OrderCardProps) {
  const [details, setDetails] = useState(false);
  const [pay, setPay] = useState(true);
  const [showPayModal, setShowPayModal] = useState<boolean>(false);
  const [showCancelModal, setShowCancelModal] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);
  const [image, setImage] = useState<File>();
  const [imageErr, setImageErr] = useState(false);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const {
    error: errPay,
    isLoading: isLoadingPay,
    fetchData: fetchPay,
  } = useAxios<Response<undefined>>(
    API_USER_ORDER_DETAILS +
      "/" +
      order.id +
      API_USER_ORDER_DETAILS_SUFF_PAYMENT,
    API_METHOD_PATCH,
    true
  );

  const {
    error: errCancel,
    isLoading: isLoadingCancel,
    fetchData: fetchCancel,
  } = useAxios<Response<undefined>>(
    API_USER_ORDERS + "/" + order.id,
    API_METHOD_DELETE
  );

  function handleImageUploadClick() {
    fileInputRef.current?.click();
  }

  function updateImage(event: React.ChangeEvent<HTMLInputElement>) {
    if (event.target.files && event.target.files[0]) {
      const file = event.target.files[0];

      const extension = file.name.split(".").pop();
      if (!(extension && types.includes(extension))) {
        setImageErr(true);
        return;
      }

      setImageErr(false);
      setImage(file);
    }
  }

  function handleCancel() {
    setPay(false);
    fetchCancel();
    setShowResult(true);
    setShowCancelModal(false);
  }

  function handlePay() {
    setPay(true);
    fetchPay({ file: image });
    setImage(undefined);
    setShowResult(true);
    setShowPayModal(false);
  }

  return (
    <>
      {order.order_details[0].status.toLowerCase() === ORDER_STATUS_PENDING ? (
        <div className="w-full bg-white shadow-md p-6 md:px-14 md:py-12 border border-primary-gray/20 rounded-lg flex flex-col gap-4 md:gap-10">
          <div className="flex flex-col gap-4 md:flex-row md:justify-between">
            <div className="flex flex-col gap-2.5">
              <p className="text-primary-black text-lg font-semibold">
                Order Date
              </p>
              <p className="text-primary-green font-bold">
                {parseTimestamp(order.order_date)}
              </p>
            </div>
            <div className="flex flex-col gap-2.5">
              <p className="text-primary-black text-lg font-semibold">
                Total Price
              </p>
              <p className="text-primary-green font-bold">
                {order.total_price ? "Rp" + order.total_price : ""}
              </p>
            </div>
            <div className="flex flex-col gap-2.5">
              <p className="text-primary-black text-lg font-semibold">
                Order Status
              </p>
              <p
                className={
                  "font-bold" + " " + statusColor(ORDER_STATUS_WAITING)
                }
              >
                {CapitalizeFirstLetter(ORDER_STATUS_WAITING)}
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
            <div className="flex flex-col md:flex-row w-full md:w-fit gap-2.5 justify-center flex-wrap">
              <div>
                <Button
                  onClick={() => {
                    setShowPayModal(true);
                  }}
                >
                  <div className="px-4 flex gap-2 justify-center items-center">
                    <Camera />
                    <p>Upload payment</p>
                  </div>
                </Button>
              </div>
              <div>
                <Button
                  type="ghost-red"
                  onClick={() => {
                    setShowCancelModal(true);
                  }}
                >
                  <div className="px-4">
                    <p>
                      Cancel<span className="hidden md:inline"> order</span>
                    </p>
                  </div>
                </Button>
              </div>
            </div>
          </div>
          {details && (
            <div className="flex flex-col gap-2.5">
              {order.order_details.map((item) => {
                return item.order_products.map((product) => {
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
                          <div className="flex gap-1.5 items-center text-primary-gray ">
                            <Hospital className="h-5" />
                            <p className="text-base truncate">
                              {item.pharmacy_name}
                            </p>
                          </div>
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
                });
              })}
            </div>
          )}
        </div>
      ) : (
        order.order_details.map((item) => {
          return (
            <OrderCardPaid
              key={item.details_id}
              order_date={order.order_date}
              order={item}
              refetchData={refetchData}
            />
          );
        })
      )}
      {showPayModal &&
        createPortal(
          <Modal
            onClose={() => {
              setImage(undefined);
              setShowPayModal(false);
            }}
          >
            <div className="w-full py-8 px-6 flex flex-col items-center justify-center gap-6">
              <h1 className="px-2 font-semibold text-primary-black text-2xl md:text-3xl text-center">
                Upload payment proof
              </h1>
              {imageErr && (
                <p className="text-primary-red text-center">
                  Allowed file types are {types.join(", ")}.
                </p>
              )}
              <div
                onClick={handleImageUploadClick}
                className="max-w-full border border-dashed border-primary-gray/60 cursor-pointer rounded-lg px-20 py-10"
              >
                {image ? (
                  <img
                    src={URL.createObjectURL(image)}
                    alt="uploaded-img"
                    className="max-h-[50vh]"
                  />
                ) : (
                  <>
                    <Image className="w-[100px] h-[100px] md:w-[130px] md:h-[130px] mx-auto text-primary-green" />
                    <p className="text-primary-gray">Max. file size 1 MB</p>
                  </>
                )}
                <input
                  ref={fileInputRef}
                  className="hidden w-full"
                  type="file"
                  onChange={updateImage}
                />
              </div>
              <div className="w-full flex flex-wrap gap-3 md:gap-7 justify-center items-center">
                <div className="w-[120px] md:w-[163px]">
                  <Button
                    submit={false}
                    size="md"
                    onClick={handlePay}
                    disabled={!image}
                  >
                    Yes
                  </Button>
                </div>
                <div className="w-[120px] md:w-[163px]">
                  <Button
                    submit={false}
                    type="ghost"
                    size="md"
                    onClick={() => {
                      setImage(undefined);
                      setShowPayModal(false);
                    }}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            </div>
          </Modal>,
          document.body
        )}
      {showCancelModal &&
        createPortal(
          <ConfirmBox
            type="cancel"
            onYes={handleCancel}
            onCancel={() => setShowCancelModal(false)}
          />,
          document.body
        )}
      {showResult &&
        createPortal(
          isLoadingPay || isLoadingCancel ? (
            <LoaderBox />
          ) : errPay || errCancel ? (
            <ErrorBox onClose={() => setShowResult(false)}>
              <p className="font-bold text-primary-black text-4xl">Oops...</p>
              <p className="text-primary-black text-xl">
                {(errPay && errPay[0].field === "server") ||
                (errCancel && errCancel[0].field === "server")
                  ? "Something is wrong, please try again"
                  : (errPay && CapitalizeFirstLetter(errPay[0].detail)) ||
                    (errCancel && CapitalizeFirstLetter(errCancel[0].detail))}
              </p>
            </ErrorBox>
          ) : (
            <SuccessBox
              onClose={() => {
                setShowResult(false);
                refetchData();
              }}
            >
              {pay ? "Payment added" : "Order cancelled"} successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </>
  );
}
