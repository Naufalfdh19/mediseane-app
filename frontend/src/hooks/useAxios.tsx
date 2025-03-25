import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

import axiosInstance from "../api/axiosInstances";
import { RemoveAccessToken } from "../utils/localstorage";
import { Err } from "../types/response";
import { API_METHOD_GET, PATH_LOGIN } from "../const/const";

function useAxios<TData, TBody = unknown>(
  endpoint: string,
  method = API_METHOD_GET,
  isMultiplePart = false,
  lazy = false,
  cb?: (data: TData) => void
) {
  const immediate = !lazy && method === API_METHOD_GET;

  const [data, setData] = useState<TData>();
  const [error, setError] = useState<Err[]>();
  const [isLoading, setIsLoading] = useState(immediate);

  const navigate = useNavigate();

  const fetchData = useCallback(
    async (body?: TBody) => {
      setIsLoading(true);
      setError(undefined);

      try {
        let header = {
          "Content-Type": "application/json",
        };

        if (isMultiplePart) {
          header = {
            "Content-Type": "multipart/form-data",
          };
        }

        const response = await axiosInstance({
          url: endpoint,
          method: method,
          headers: header,
          data: body,
        });

        setData(response.data);

        if (response) {
          cb?.(response.data);
        }

        return response.data as TData;
      } catch (err) {
        if (axios.isAxiosError(err)) {
          if (err.response) {
            if (err.response.status === 401) {
              RemoveAccessToken();
              navigate(PATH_LOGIN);
            }
            setError(err.response.data.error);
          } else if (err.request) {
            setError([{ field: "server", detail: "internal server error" }]);
          }
        } else {
          setError([{ field: "server", detail: "internal server error" }]);
        }
      } finally {
        setIsLoading(false);
      }
    },
    [endpoint, JSON.stringify(method)]
  );

  useEffect(() => {
    if (immediate) {
      fetchData();
    }
  }, [fetchData, immediate]);

  return { data, isLoading, error, fetchData };
}

export default useAxios;
