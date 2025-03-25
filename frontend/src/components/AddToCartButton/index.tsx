import { useEffect, useState } from "react";
import { ShoppingCart } from "lucide-react";
import { createPortal } from "react-dom";

import Button from "../Button";
import ConfirmBox from "../ConfirmBox";
import ErrorBox from "../ErrorBox";

import {
  useCartState,
  useRemoveCartItem,
  useSetCartQuantity,
} from "../../store/cart";
import useAxios from "../../hooks/useAxios";
import { Response } from "../../types/response";
import {
  API_METHOD_DELETE,
  API_METHOD_PUT,
  API_USER_CARTS,
} from "../../const/const";
import { useCart } from "../../hooks/useCart";
import { CapitalizeFirstLetter } from "../../utils/formatter";

interface AddToCartButtonProps {
  product_id: number;
  size?: "sm" | "md";
}

export default function AddToCartButton({
  product_id,
  size = "md",
}: AddToCartButtonProps) {
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<boolean>(false);
  const [showError, setShowError] = useState<boolean>(false);

  const cartState = useCartState();
  const setCartQuantity = useSetCartQuantity();
  const removeCartItem = useRemoveCartItem();
  const { refetchCart } = useCart();

  const cart = cartState.items.find(
    (item) => item.pharmacy_product_id === product_id
  );

  const {
    fetchData: fetchUpdate,
    isLoading: isUpdateLoading,
    error: UpdateErr,
  } = useAxios<Response<undefined>>(API_USER_CARTS, API_METHOD_PUT);

  const { fetchData: fetchDelete, error: DeleteErr } = useAxios<
    Response<undefined>
  >(API_USER_CARTS + (cart ? "/" + cart?.cart_item_id : ""), API_METHOD_DELETE);

  function handleDecrease() {
    if (cart) {
      if (cart.quantity > 1) {
        fetchUpdate({
          pharmacy_product_id: product_id,
          quantity: -1,
        }).then((res) => {
          if (res && res.message) {
            setCartQuantity(product_id, cart.quantity - 1);
          }
        });
      } else {
        setShowDeleteConfirm(true);
      }
    }
  }

  function handleIncrease() {
    if (cart) {
      fetchUpdate({
        pharmacy_product_id: product_id,
        quantity: 1,
      }).then((res) => {
        if (res && res.message) {
          setCartQuantity(product_id, cart.quantity + 1);
        }
      });
    }
  }

  function handleDelete() {
    fetchDelete().then((res) => {
      if (res && res.message) {
        removeCartItem(product_id);
      }
    });

    setShowDeleteConfirm(false);
  }

  function handleAddToCart() {
    fetchUpdate({
      pharmacy_product_id: product_id,
      quantity: 1,
    }).finally(() => refetchCart());
  }

  useEffect(() => {
    if (UpdateErr || DeleteErr) {
      setShowError(true);
    }
  }, [UpdateErr, DeleteErr]);

  return (
    <>
      {cart ? (
        <div className="w-fit flex gap-3 items-center self-center">
          <div
            className={`${size === "md" ? "w-[42px]" : "w-[28px]"} ${
              isUpdateLoading && "animate-pulse bg-secondary-gray/20 rounded"
            }`}
            onClick={(e) => e.stopPropagation()}
          >
            <Button
              type="ghost-green"
              size={size === "md" ? "sm" : "xs"}
              square={true}
              onClick={handleDecrease}
              disabled={isUpdateLoading}
            >
              -
            </Button>
          </div>
          <div className="w-[30px]">
            <p className="text-center text-primary-black">{cart.quantity}</p>
          </div>
          <div
            className={`${size === "md" ? "w-[42px]" : "w-[28px]"} ${
              isUpdateLoading && "animate-pulse bg-secondary-gray/20 rounded"
            }`}
            onClick={(e) => e.stopPropagation()}
          >
            <Button
              type="ghost-green"
              size={size === "md" ? "sm" : "xs"}
              square={true}
              disabled={isUpdateLoading}
              onClick={handleIncrease}
            >
              +
            </Button>
          </div>
        </div>
      ) : (
        <div className="w-full" onClick={(e) => e.stopPropagation()}>
          <Button onClick={handleAddToCart} size={size}>
            <div className="flex justify-center items-center gap-1.5">
              <ShoppingCart className="hidden md:inline" />
              <p>Add to cart</p>
            </div>
          </Button>
        </div>
      )}
      {showDeleteConfirm &&
        createPortal(
          <ConfirmBox
            type="delete"
            onYes={() => {
              handleDelete();
            }}
            onCancel={() => {
              setShowDeleteConfirm(false);
            }}
          />,
          document.body
        )}
      {showError &&
        createPortal(
          <ErrorBox onClose={() => setShowError(false)}>
            <p className="font-bold text-primary-black text-4xl">Oops...</p>
            <p className="text-primary-black text-xl">
              {(UpdateErr && UpdateErr[0].field === "server") ||
              (DeleteErr && DeleteErr[0].field === "server")
                ? "Something is wrong, please try again"
                : (UpdateErr && CapitalizeFirstLetter(UpdateErr[0].detail)) ||
                  (DeleteErr && CapitalizeFirstLetter(DeleteErr[0].detail))}
            </p>
          </ErrorBox>,
          document.body
        )}
    </>
  );
}
