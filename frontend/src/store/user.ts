import { useDispatch, useSelector } from "react-redux";

import { AppState } from "./store";
import {
  setUserState,
  UserState,
  resetUserState,
  setUserUsername,
  setLoading,
} from "./userSlice";

export const useUserState = () => useSelector((state: AppState) => state.user);

export const useSetUserState = () => {
  const dispatch = useDispatch();

  return (userData: UserState) => {
    dispatch(setUserState(userData));
  };
};

export const useResetUserState = () => {
  const dispatch = useDispatch();
  return () => dispatch(resetUserState());
};

export const useSetUserUsername = () => {
  const dispatch = useDispatch();
  return (username: string) => dispatch(setUserUsername(username));
};

export const useSetLoading = () => {
  const dispatch = useDispatch();
  return () => dispatch(setLoading());
};
