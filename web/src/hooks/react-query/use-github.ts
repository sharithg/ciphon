import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { API_URL } from "./constants";
import { toast } from "../use-toast";
import { fetchData } from ".";

export type TConnectRepo = {
  name: string;
  owner: string;
};

export type TGithubRepo = {
  repoId: number;
  name: string;
  owner: string;
  description: string;
  url: string;
  repoCreatedAt: string;
};

export type TNewGithubRepos = {
  id: number;
  name: string;
  description: string;
  lastUpdated: string;
  owner: string;
};

export const useGetRepos = () => {
  return useQuery({
    queryKey: ["repos"],
    queryFn: () => fetchData<{ data: TGithubRepo[] }>(`${API_URL}/repos`),
  });
};

export const useGetNewRepos = () => {
  return useQuery({
    queryKey: ["new-repos"],
    queryFn: () => fetchData<{ data: TGithubRepo[] }>(`${API_URL}/new-repos`),
  });
};

export const useConnectRepo = () => {
  const mutation = useMutation({
    mutationFn: (data: TConnectRepo) => {
      return axios.post(`${API_URL}/repos/connect`, data, {
        headers: {
          "Content-Type": "application/json",
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
