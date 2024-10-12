import axios from "axios";
import { refreshAccessToken } from "../hooks/user-auth";

export const apiClient = axios.create();

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken");
  if (token) {
    config.headers["Authorization"] = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem("refreshToken");
        if (!refreshToken) {
          window.location.href = "/login";
          return;
        }
        const response = await refreshAccessToken(refreshToken);

        localStorage.setItem("accessToken", response.accessToken);
        apiClient.defaults.headers["Authorization"] =
          `Bearer ${response.accessToken}`;
        originalRequest.headers["Authorization"] =
          `Bearer ${response.accessToken}`;

        return apiClient(originalRequest);
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
      } catch (e) {
        console.log("Refresh token expired, logging out.");
        window.location.href = "/login";
      }
    }

    return Promise.reject(error);
  }
);
