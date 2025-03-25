import { useDispatch, useSelector } from "react-redux";
import { AppState } from "./store";
import { resetCheckout, toggleCheckout } from "./checkoutSlice";

export const useCheckoutState = () =>
  useSelector((state: AppState) => state.checkout.cartIDs);

const selectIsCartInCheckout = (cartID: number) => (state: AppState) =>
  state.checkout.cartIDs.includes(cartID);

export const useIsCartInCheckout = (cartID: number) =>
  useSelector((state: AppState) => selectIsCartInCheckout(cartID)(state));

export const useToggleCheckout = () => {
  const dispatch = useDispatch();

  return (cartID: number) => {
    dispatch(toggleCheckout(cartID));
  };
};

export const useResetCheckout = () => {
  const dispatch = useDispatch();

  return () => {
    dispatch(resetCheckout());
  };
};
