import { GithubRepo } from "@/@types/api";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export const useGetRepos = () => {
  return useQuery({
    queryKey: ["todos"],
    queryFn: () => axios.get<GithubRepo[]>("http://localhost:8000/repos"),
  });
};
