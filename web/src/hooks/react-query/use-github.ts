import { useMutation, useQuery } from "@tanstack/react-query";
import { API_URL } from "./constants";
import { toast } from "../use-toast";
import { fetchData } from ".";
import { withJwt } from "../user-auth";
import { apiClient } from "../../axios";
import {
  TGetAllReposRow,
  TGithubRepoResponse,
  TConnectRepoRequest,
} from "../../types/api";

export const useGetRepos = () => {
  return useQuery({
    queryKey: ["repos"],
    queryFn: () => fetchData<TGetAllReposRow[]>(`${API_URL}/repos`),
  });
};

export const useGetNewRepos = () => {
  return useQuery({
    queryKey: ["new-repos"],
    queryFn: () => fetchData<TGithubRepoResponse[]>(`${API_URL}/repos/new`),
  });
};

export const useConnectRepo = () => {
  const mutation = useMutation({
    mutationFn: (data: TConnectRepoRequest) => {
      return apiClient.post(`${API_URL}/repos/connect`, data, {
        headers: {
          "Content-Type": "application/json",
          ...withJwt(),
        },
      });
    },
    onSuccess: () => {
      toast({
        title: "Succesfully connected repo",
      });
    },
    onError: () => {
      toast({
        title: "Error connecting repo",
        variant: "destructive",
      });
    },
  });
  return mutation;
};
