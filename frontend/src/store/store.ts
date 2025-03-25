import { combineReducers, configureStore } from "@reduxjs/toolkit";
import userSlice from "./userSlice";
import cartSlice from "./cartSlice";
import checkoutSlice from "./checkoutSlice";

const rootReducer = combineReducers({
  user: userSlice.reducer,
  cart: cartSlice.reducer,
  checkout: checkoutSlice.reducer,
});

export const setupStore = (preloadedState?: Partial<AppState>) => {
  return configureStore({
    reducer: rootReducer,
    preloadedState,
  });
};

export type AppState = ReturnType<typeof rootReducer>;
export type AppStore = ReturnType<typeof setupStore>;
export type AppDispatch = AppStore["dispatch"];
