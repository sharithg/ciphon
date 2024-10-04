import axios, { AxiosResponse } from "axios";

export const fetchData = async <T>(url: string): Promise<T> => {
  const response: AxiosResponse<T> = await axios.get(url);
  return response.data;
};
