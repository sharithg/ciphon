import { API_URL } from "./constants";
import { useMutation, useQuery } from "@tanstack/react-query";
import { fetchData } from ".";
import axios from "axios";

type WorklfowRun = {
  commitSha: string;
  repoName: string;
  pipelineId: string;
  workflowId: string;
  workflowName: string;
  status: string;
  branch: string;
  createdAt: string;
  duration: number | null;
};

export const useGetWorkflows = () => {
  return useQuery({
    queryKey: ["nodes"],
    queryFn: () => fetchData<{ data: WorklfowRun[] }>(`${API_URL}/workflows`),
  });
};

export const useRunWorkflow = () => {
  const mutation = useMutation({
    mutationFn: (workflowId: string) => {
      return axios.post(`${API_URL}/workflows/trigger/${workflowId}`);
    },
  });
  return mutation;
};
