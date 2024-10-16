import { API_URL } from "./react-query/constants";
import { fetchData, WithData } from "./react-query";
import { useAtom } from "jotai";
import { userAtom } from "../components/atoms/user";
import { apiClient } from "../axios";
import { TGetUserByIdRow, TTokenPair } from "../types/api";

function isTokenExpired(token: string) {
  const base64Url = token.split(".")[1];
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map(function (c) {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );

  const { exp } = JSON.parse(jsonPayload);

  const currentTime = Math.floor(Date.now() / 1000);
  return exp < currentTime;
}

export const withJwt = () => {
  return {
    Authorization: `Bearer ${localStorage.getItem("accessToken")}`,
  };
};

const getToken = async (code: string) => {
  const tok = await fetchData<TTokenPair>(
    `${API_URL}/auth/login/github/callback?code=${code}`
  );
  return tok;
};

const getUser = async () => {
  const tok = await fetchData<TGetUserByIdRow>(`${API_URL}/user`);
  return tok;
};

export const refreshAccessToken = async (token: string) => {
  const tok = await apiClient.post<WithData<TTokenPair>>(
    `${API_URL}/auth/refresh-token`,
    { token },
    {
      headers: {
        ...withJwt(),
      },
    }
  );
  return tok.data.data;
};

export const isAuthenticated = async () => {
  const token = localStorage.getItem("accessToken");
  if (!token) {
    return false;
  }
  const refreshToken = localStorage.getItem("refreshToken");
  if (!refreshToken) {
    return false;
  }

  if (isTokenExpired(token)) {
    const newTokens = await refreshAccessToken(refreshToken);
    localStorage.setItem("accessToken", newTokens.accessToken);
    localStorage.setItem("refreshToken", newTokens.refreshToken);
  }
  const user = await getUser();
  return !!user;
};

export const useAuth = () => {
  const [user, setUser] = useAtom(userAtom);

  const handleCallback = async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get("code");
    if (!code) {
      throw new Error("code not found");
    }
    const token = await getToken(code);
    localStorage.setItem("accessToken", token.accessToken);
    localStorage.setItem("refreshToken", token.refreshToken);

    const user = await getUser();

    setUser(user);
  };

  return {
    handleCallback,
    user,
  };
};
