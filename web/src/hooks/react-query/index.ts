import { AxiosResponse } from "axios";
import { withJwt } from "../user-auth";
import { apiClient } from "../../axios";

export const fetchData = async <T>(url: string): Promise<T> => {
  const response: AxiosResponse<WithData<T>> = await apiClient.get(url, {
    headers: {
      ...withJwt(),
    },
  });
  return response.data.data;
};

export type WithData<T> = { data: T };
