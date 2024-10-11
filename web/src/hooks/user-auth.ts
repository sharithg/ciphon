import axios from "axios";
import { API_URL } from "./react-query/constants";

const getToken = async (code: string) => {
  const tok = await axios.get<{ data: { token: string } }>(
    `${API_URL}/auth/login/github/callback?code=${code}`
  );
  return tok.data.data.token;
};

const useAuth = () => {};
