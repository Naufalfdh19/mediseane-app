import { Dispatch, SetStateAction } from "react";
import { ChevronDown, Truck } from "lucide-react";

import CheckoutItem from "./CheckoutItem";
import ErrorCard from "../../../components/ErrorCard";

import useAxios from "../../../hooks/useAxios";
import {
  Delivery,
  DeliveryData,
  GroupedCartItem,
  Response,
} from "../../../types/response";
import { API_DELIVERY, API_METHOD_GET } from "../../../const/const";

interface GroupedCheckoutCardProps {
  item: GroupedCartItem;
  isAddressSet: boolean;
  shippingID: number | undefined;
  deliveries: DeliveryData[] | undefined;
  setDeliveries: Dispatch<SetStateAction<DeliveryData[] | undefined>>;
  deliveryCosts: Map<number, number>;
  setDeliveryCosts: Dispatch<SetStateAction<Map<number, number>>>;
}

export default function GroupedCheckoutCard({
  item,
  isAddressSet,
  shippingID,
  deliveries,
  setDeliveries,
  deliveryCosts,
  setDeliveryCosts,
}: GroupedCheckoutCardProps) {
  const { data, error, isLoading, fetchData } = useAxios<Response<Delivery[]>>(
    API_DELIVERY + item.pharmacy_id,
    API_METHOD_GET,
    false,
    true
  );

  const shippingCost = data?.data.find((item) => item.id === shippingID)?.cost;
  const ShippingEtd = data?.data.find((item) => item.id === shippingID)?.etd;

  function handleFetchShipping() {
    if (data === undefined && isAddressSet) fetchData();
  }

  function handleShippingChange(e: React.ChangeEvent<HTMLSelectElement>) {
    if (deliveries) {
      const newDeliveries: DeliveryData[] = deliveries.map((delivery) => {
        if (delivery.pharmacy_id === item.pharmacy_id) {
          delivery.delivery_id = +e.target.value;
          return delivery;
        } else {
          return delivery;
        }
      });
      setDeliveries(newDeliveries);
      const newMap = deliveryCosts;
      newMap.set(
        item.pharmacy_id,
        data!.data.find((item) => item.id === +e.target.value)!.cost
      );
      setDeliveryCosts(newMap);
    }
  }

  function totalItemPrice() {
    let total = 0;
    item.items.map((item) => (total += +item.subtotal));
    return total;
  }

  return (
    <div className="w-full bg-white shadow-md p-4 border border-primary-gray/20 rounded-lg flex flex-col gap-2">
      <div className="text-primary-black font-semibold overflow-auto">
        {item.pharmacy_name}
      </div>
      {item.items.map((item) => (
        <CheckoutItem key={item.cart_item_id} item={item} />
      ))}
      {error ? (
        <ErrorCard errors={error} />
      ) : isLoading ? (
        <div className="w-full h-[150px] bg-white p-4 border-y border-primary-gray/50 flex flex-col gap-2 overflow-auto">
          <p>Shipping</p>
          <div className="w-full px-[13px] py-[12px] h-[50px] rounded-lg bg-primary-gray/20 animate-pulse appearance-none focus:outline-none"></div>
        </div>
      ) : data ? (
        <div className="w-full h-[150px] bg-white p-4 border-y border-primary-gray/50 flex flex-col gap-2 overflow-auto">
          <p>Shipping</p>
          <div className="relative">
            <select
              className="w-full px-[13px] py-[12px] rounded-lg bg-white border border-primary-black appearance-none focus:outline-none"
              onChange={handleShippingChange}
            >
              <option value="" disabled selected hidden>
                Choose shipping method
              </option>
              {data.data &&
                data.data.map((options) => (
                  <option key={options.id} value={options.id}>
                    {options.name}
                  </option>
                ))}
            </select>
            <ChevronDown className="absolute right-3 top-[50%] translate-y-[-50%] text-primary-black" />
          </div>
          <div className="flex gap-2 text-primary-green overflow-auto">
            <Truck strokeWidth={1.4} className="flex-shrink-0" />
            <p>Guaranteed arrival in {ShippingEtd ? ShippingEtd : "---"}</p>
          </div>
        </div>
      ) : (
        <div
          className={`w-full h-[150px] ${
            isAddressSet
              ? "bg-primary-green bg-opacity-[18%] text-primary-green cursor-pointer"
              : "bg-secondary-gray text-primary-gray cursor-not-allowed"
          } p-4  rounded-lg flex  gap-2 items-center justify-center overflow-auto`}
          onClick={handleFetchShipping}
        >
          {isAddressSet ? (
            <p>Click here to calculate shipping cost</p>
          ) : (
            <p>
              Set an address first before calculating{" "}
              <span className="font-bold">shipping cost</span>
            </p>
          )}
          <Truck strokeWidth={1.4} className="flex-shrink-0" />
        </div>
      )}
      <div className="w-full bg-white p-4 b rounded-lg flex flex-col gap-2 overflow-auto">
        <div className="flex justify-between gap-2">
          <p className="text-primary-black">Subtotal: </p>
          <p className="text-primary-green">Rp{totalItemPrice()}</p>
        </div>

        <div className="flex justify-between gap-2">
          <p className="text-primary-black">Shipping cost: </p>
          <p className="text-primary-green">
            Rp{shippingCost ? shippingCost : "---"}
          </p>
        </div>
        <div className="flex justify-between gap-2">
          <p className="text-primary-black">Total: </p>
          <p className="text-primary-green font-semibold">
            Rp{totalItemPrice() + (shippingCost ? +shippingCost : 0)}
          </p>
        </div>
      </div>
    </div>
  );
}
