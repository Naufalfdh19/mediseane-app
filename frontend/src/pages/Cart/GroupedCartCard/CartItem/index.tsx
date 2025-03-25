import { useState } from "react";
import { createPortal } from "react-dom";
import { Check } from "lucide-react";

import ErrorBox from "../../../../components/ErrorBox";
import Button from "../../../../components/Button";
import Modal from "../../../../components/Modal/Modal";
import PageLoader from "../../../../components/PageLoader";
import ConfirmBox from "../../../../components/ConfirmBox";

import useAxios from "../../../../hooks/useAxios";
import { CapitalizeFirstLetter } from "../../../../utils/formatter";
import {
  useIsCartInCheckout,
  useToggleCheckout,
} from "../../../../store/checkout";
import { CartItem } from "../../../../types/response";
import {
  API_CART,
  API_METHOD_DELETE,
  API_METHOD_PUT,
} from "../../../../const/const";
import { useCart } from "../../../../hooks/useCart";

interface CartItemProps {
  item: CartItem;
  refetch: () => void;
}

export default function CartItemBox({ item, refetch }: CartItemProps) {
  const [showResult, setShowResult] = useState<boolean>(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<boolean>(false);
  const toggleCheckout = useToggleCheckout();
  const isInCheckout = useIsCartInCheckout(item.cart_item_id);
  const { refetchCart } = useCart();

  const {
    error: errorUpdate,
    isLoading: isUpdateLoading,
    fetchData: updateCart,
  } = useAxios(API_CART, API_METHOD_PUT);

  const {
    error: errorDelete,
    isLoading: isDeleteLoading,
    fetchData: deleteCart,
  } = useAxios(API_CART + "/" + item.cart_item_id, API_METHOD_DELETE);

  function handlePlusButtonClick() {
    updateCart({
      pharmacy_product_id: item.pharmacy_product_id,
      quantity: 1,
    }).finally(() => {
      if (errorUpdate) {
        setShowResult(true);
      }
      refetch();
    });
  }

  function handleDeleteCart() {
    deleteCart().finally(() => {
      if (!errorDelete) {
        setShowResult(false);
        if (isInCheckout) {
          toggleCheckout(item.cart_item_id);
        }
      }
      refetchCart();
      refetch();
    });
    setShowResult(true);
  }

  function handleMinButtonClick() {
    if (item.quantity > 1) {
      updateCart({
        pharmacy_product_id: item.pharmacy_product_id,
        quantity: -1,
      }).finally(() => {
        refetch();
      });
    } else {
      setShowDeleteConfirm(true);
    }
  }

  return (
    <div
      className="h-[150px] border-t-[1px] border-primary-gray/50 py-3 flex items-center overflow-auto gap-3 md:gap-6 cursor-pointer"
      onClick={() => toggleCheckout(item.cart_item_id)}
    >
      <div className="ml-4">
        {isInCheckout ? (
          <div className="bg-primary-green h-[20px] w-[20px] rounded flex items-center justify-center">
            <Check
              className="h-[16px] w-[16px] text-primary-white"
              strokeWidth={3.5}
            />
          </div>
        ) : (
          <div className="h-[20px] w-[20px] border-[2px] rounded border-primary-gray/50 flex-shrink-0"></div>
        )}
      </div>

      <div className="flex gap-3 items-center">
        <div className="w-[60px] h-[60px] md:w-[80px] md:h-[80px] flex-shrink-0`">
          <img src={item.image} alt="image" />
        </div>
        <div className="grow flex flex-col gap-2">
          <div className="w-full max-w-[5000px]">
            <p className="truncate text-primary-black">{item.name}</p>
            <p className="truncate text-sm text-primary-gray">
              {item.manufacturer}
            </p>
          </div>
          <p className="truncate text-primary-green font-semibold">
            Rp{item.subtotal}
          </p>
          <div className="w-fit flex gap-3 items-center">
            <div
              className={`w-[28px] ${
                isUpdateLoading && "animate-pulse bg-secondary-gray/20 rounded"
              }`}
              onClick={(e) => e.stopPropagation()}
            >
              <Button
                type="ghost-green"
                size="xs"
                square={true}
                onClick={handleMinButtonClick}
                disabled={isUpdateLoading}
              >
                -
              </Button>
            </div>
            <div className="w-[30px]">
              <p className="text-center text-primary-black">{item.quantity}</p>
            </div>
            <div
              className={`w-[28px] ${
                isUpdateLoading && "animate-pulse bg-secondary-gray/20 rounded"
              }`}
              onClick={(e) => e.stopPropagation()}
            >
              <Button
                type="ghost-green"
                size="xs"
                square={true}
                onClick={handlePlusButtonClick}
                disabled={isUpdateLoading}
              >
                +
              </Button>
            </div>
          </div>
        </div>
      </div>
      {showDeleteConfirm &&
        createPortal(
          <ConfirmBox
            type="delete"
            onYes={handleDeleteCart}
            onCancel={() => setShowDeleteConfirm(false)}
          />,
          document.body
        )}
      {showResult &&
        createPortal(
          isDeleteLoading ? (
            <Modal onClose={() => setShowResult(false)}>
              <PageLoader />
            </Modal>
          ) : (
            (errorUpdate || errorDelete) && (
              <ErrorBox onClose={() => setShowResult(false)}>
                <p className="font-bold text-primary-black text-4xl">Oops...</p>
                <p className="text-primary-black text-xl">
                  {(errorUpdate && errorUpdate[0].field === "server") ||
                  (errorDelete && errorDelete[0].field === "server")
                    ? "Something is wrong, please try again"
                    : (errorUpdate &&
                        CapitalizeFirstLetter(errorUpdate[0].detail)) ||
                      (errorDelete &&
                        CapitalizeFirstLetter(errorDelete[0].detail))}
                </p>
              </ErrorBox>
            )
          ),
          document.body
        )}
    </div>
  );
}
