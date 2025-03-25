import axios from "axios";
import config from "../config";
import { GetAccessToken } from "../utils/localstorage";

const axiosInstance = axios.create({
  baseURL: config.API_BASE_URL,
  timeout: 600000,
});

axiosInstance.interceptors.request.use(
  (config) => {
    const token = GetAccessToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default axiosInstance;
