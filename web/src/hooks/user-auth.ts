import axios from "axios";
import { API_URL } from "./react-query/constants";
import { WithData } from "./react-query";
import { useAtom } from "jotai";
import { userAtom } from "../components/atoms/user";

export type User = {
  id: string;
  username: string;
  email: string;
  avatarUrl: string;
};

const withJwt = () => {
  return {
    Authorization: `Bearer ${localStorage.getItem("jwtToken")}`,
  };
};

const getToken = async (code: string) => {
  const tok = await axios.get<{ data: { token: string } }>(
    `${API_URL}/auth/login/github/callback?code=${code}`
  );
  return tok.data.data.token;
};

const getUser = async () => {
  const tok = await axios.get<WithData<User>>(`${API_URL}/user`, {
    headers: {
      ...withJwt(),
    },
  });
  return tok.data.data;
};

export const isAuthenticated = async () => {
  const token = localStorage.getItem("jwtToken");
  if (!token) {
    return false;
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
    localStorage.setItem("jwtToken", token);

    const user = await getUser();
    setUser(user);
  };

  return {
    handleCallback,
    user,
  };
};
