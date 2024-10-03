import { GithubRepo } from "@/@types/api";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { API_URL } from "./constants";

export const useGetRepos = () => {
  return useQuery({
    queryKey: ["todos"],
    queryFn: () => axios.get<GithubRepo[]>(`${API_URL}/repos`),
  });
};
