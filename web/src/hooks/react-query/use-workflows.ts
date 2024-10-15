import { API_URL } from "./constants";
import { useMutation, useQuery } from "@tanstack/react-query";
import { fetchData, WithData } from ".";
import { useAtom } from "jotai";
import { jobs, workflows } from "../../components/atoms/workflows";
import { withJwt } from "../user-auth";
import { apiClient } from "../../axios";
import {
  TCommandOutput,
  TJobs,
  TSteps,
  TWorkflowRunInfo,
} from "../../types/api";

export const useGetWorkflows = () => {
  const [, setWorkflows] = useAtom(workflows);

  return useQuery({
    queryKey: ["workflows"],
    queryFn: () =>
      fetchData<WithData<TWorkflowRunInfo[]>>(`${API_URL}/workflows`),
    onSuccess: (data) => {
      setWorkflows(data.data);
    },
  });
};

export const useRunWorkflow = () => {
  const mutation = useMutation({
    mutationFn: (workflowId: string) => {
      return apiClient.post(
        `${API_URL}/workflows/trigger/${workflowId}`,
        {},
        {
          headers: {
            ...withJwt(),
          },
        }
      );
    },
  });
  return mutation;
};

export const useGetJobs = (workflowId: string) => {
  const [, setJobs] = useAtom(jobs);

  return useQuery({
    queryKey: [`workflows/${workflowId}/jobs`],
    queryFn: () =>
      fetchData<TJobs[]>(`${API_URL}/workflows/${workflowId}/jobs`),
    onSuccess: (data) => {
      setJobs(data);
    },
  });
};

export const useGetSteps = (workflowId: string, jobId: string) => {
  return useQuery({
    queryKey: [`workflows/${workflowId}/jobs`],
    queryFn: () =>
      fetchData<WithData<TSteps[]>>(
        `${API_URL}/workflows/${workflowId}/jobs/${jobId}/steps`
      ),
    refetchInterval: 500,
  });
};

export const useGetCommandOutput = (
  workflowId: string,
  jobId: string,
  stepId: string,
  enabled = true
) => {
  return useQuery({
    queryKey: [`workflows/${workflowId}/jobs/${jobId}/steps/${stepId}/output`],
    queryFn: () =>
      fetchData<TCommandOutput[]>(
        `${API_URL}/workflows/${workflowId}/jobs/${jobId}/steps/${stepId}/output`
      ),
    refetchInterval: 500,
    enabled,
  });
};
