import { useDispatch, useSelector } from "react-redux";

import { AppState } from "./store";
import {
  addCartItem,
  removeCartItem,
  resetCartState,
  setCartQuantity,
  setCartState,
} from "./cartSlice";
import { CartOverviewListItem } from "../types/response";

export const useCartState = () => useSelector((state: AppState) => state.cart);

export const useSetCartState = () => {
  const dispatch = useDispatch();
  return (cartItems: CartOverviewListItem[]) =>
    dispatch(setCartState(cartItems));
};

export const useAddCartItem = () => {
  const dispatch = useDispatch();

  return (cartItem: CartOverviewListItem) => {
    dispatch(addCartItem(cartItem));
  };
};

export const useRemoveCartItem = () => {
  const dispatch = useDispatch();

  return (pharmacy_product_id: number) => {
    dispatch(removeCartItem(pharmacy_product_id));
  };
};

export const useResetCartState = () => {
  const dispatch = useDispatch();
  return () => dispatch(resetCartState());
};

export const useSetCartQuantity = () => {
  const dispatch = useDispatch();
  return (pharmacy_product_id: number, quantity: number) =>
    dispatch(
      setCartQuantity({
        pharmacy_product_id: pharmacy_product_id,
        quantity: quantity,
      })
    );
};
