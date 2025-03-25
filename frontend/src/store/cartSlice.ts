import { createSlice, PayloadAction } from "@reduxjs/toolkit";

import { CartOverviewListItem } from "../types/response";

export interface CartState {
  items: CartOverviewListItem[];
}

const initialState: CartState = {
  items: [],
};

const cartSlice = createSlice({
  name: "cart",
  initialState,
  reducers: {
    setCartState: (state, action: PayloadAction<CartOverviewListItem[]>) => {
      state.items = action.payload;
    },
    addCartItem: (state, action: PayloadAction<CartOverviewListItem>) => {
      state.items = state.items.concat([action.payload]);
    },
    removeCartItem: (state, action: PayloadAction<number>) => {
      state.items = state.items.filter(
        (item) => item.pharmacy_product_id !== action.payload
      );
    },
    resetCartState: (state) => {
      state.items = [];
    },
    setCartQuantity: (
      state,
      action: PayloadAction<{ pharmacy_product_id: number; quantity: number }>
    ) => {
      state.items = state.items.map((item) => {
        if (item.pharmacy_product_id === action.payload.pharmacy_product_id) {
          return {
            cart_item_id: item.cart_item_id,
            pharmacy_product_id: item.pharmacy_product_id,
            quantity: action.payload.quantity,
          };
        }
        return item;
      });
    },
  },
});

export default cartSlice;

export const {
  setCartState,
  addCartItem,
  removeCartItem,
  resetCartState,
  setCartQuantity,
} = cartSlice.actions;
