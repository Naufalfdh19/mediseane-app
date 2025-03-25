import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface checkoutState {
  cartIDs: number[];
}

const initialState: checkoutState = {
  cartIDs: [],
};

const checkoutSlice = createSlice({
  name: "checkout",
  initialState,
  reducers: {
    toggleCheckout: (state, action: PayloadAction<number>) => {
      if (state.cartIDs.includes(action.payload)) {
        state.cartIDs = state.cartIDs.filter((id) => id !== action.payload);
      } else {
        state.cartIDs.push(action.payload);
      }
    },
    resetCheckout: (state) => {
      state.cartIDs = [];
    },
  },
});

export default checkoutSlice;

export const { toggleCheckout, resetCheckout } = checkoutSlice.actions;
