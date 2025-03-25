import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface UserState {
  id: number;
  name: string;
  email: string;
  profile_photo: string;
  role: string;
  is_verified: boolean;
  is_loading: boolean;
}

const initialState: UserState = {
  id: 0,
  name: "",
  email: "",
  profile_photo: "",
  role: "",
  is_verified: false,
  is_loading: false,
};

const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setUserState: (state, action: PayloadAction<UserState>) => {
      state.id = action.payload.id;
      state.name = action.payload.name;
      state.email = action.payload.email;
      state.profile_photo = action.payload.profile_photo;
      state.role = action.payload.role;
      state.is_verified = action.payload.is_verified;
    },
    resetUserState: (state) => {
      state.id = 0;
      state.name = "";
      state.email = "";
      state.profile_photo = "";
      state.role = "";
      state.is_verified = false;
    },
    setUserUsername: (state, action: PayloadAction<string>) => {
      state.name = action.payload;
    },
    setLoading: (state) => {
      state.is_loading = true;
    },
  },
});

export default userSlice;

export const {
  setUserState,
  resetUserState,
  setUserUsername,
  setLoading,
} = userSlice.actions;
