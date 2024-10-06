import { API_URL } from "./constants";
import { useMutation, useQuery } from "@tanstack/react-query";
import { fetchData, WithData } from ".";
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

type Job = {
  id: string;
  name: string;
  status: string;
};

type Step = {
  type: string;
  id: string;
  name: string;
  command: string;
  status: string | null;
};

export const useGetWorkflows = () => {
  return useQuery({
    queryKey: ["workflows"],
    queryFn: () => fetchData<WithData<WorklfowRun[]>>(`${API_URL}/workflows`),
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

export const useGetJobs = (workflowId: string) => {
  return useQuery({
    queryKey: [`workflows/${workflowId}/jobs`],
    queryFn: () =>
      fetchData<WithData<Job[]>>(`${API_URL}/workflows/${workflowId}/jobs`),
  });
};

export const useGetSteps = (workflowId: string, jobId: string) => {
  return useQuery({
    queryKey: [`workflows/${workflowId}/jobs`],
    queryFn: () =>
      fetchData<WithData<Step[]>>(
        `${API_URL}/workflows/${workflowId}/jobs/${jobId}/steps`
      ),
  });
};
