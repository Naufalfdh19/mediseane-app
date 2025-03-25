import { API_METHOD_GET, API_USER_CARTS_OVERVIEW } from "../const/const";
import { useSetCartState } from "../store/cart";
import { CartOverviewListItem, Response } from "../types/response";
import useAxios from "./useAxios";

export function useCart() {
  const setCartState = useSetCartState();

  const { isLoading, error, fetchData } = useAxios<
    Response<CartOverviewListItem[]>
  >(API_USER_CARTS_OVERVIEW, API_METHOD_GET, false, true);

  const refetchCart = () => {
    fetchData().then((res) => {
      if (res && res.data) {
        setCartState(res.data);
      }
    });
  };

  return { isLoading, error, refetchCart };
}
