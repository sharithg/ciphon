import { API_URL } from "./constants";
import { useMutation, useQuery } from "@tanstack/react-query";
import { fetchData, WithData } from ".";
import axios from "axios";
import { useAtom } from "jotai";
import { jobs, workflows } from "../../components/atoms/workflows";

export type WorklfowRun = {
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

export type Job = {
  id: string;
  name: string;
  status: string;
};

export type Step = {
  type: string;
  id: string;
  name: string;
  command: string;
  status: string | null;
};

type CommandOutput = {
  id: string;
  step_id: string;
  stdout: string;
  type?: string;
  created_at: string;
};

export const useGetWorkflows = () => {
  const [, setWorkflows] = useAtom(workflows);

  return useQuery({
    queryKey: ["workflows"],
    queryFn: () => fetchData<WithData<WorklfowRun[]>>(`${API_URL}/workflows`),
    onSuccess: (data) => {
      setWorkflows(data.data);
    },
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
  const [, setJobs] = useAtom(jobs);

  return useQuery({
    queryKey: [`workflows/${workflowId}/jobs`],
    queryFn: () =>
      fetchData<WithData<Job[]>>(`${API_URL}/workflows/${workflowId}/jobs`),
    onSuccess: (data) => {
      setJobs(data.data);
    },
  });
};

export const useGetSteps = (workflowId: string, jobId: string) => {
  return useQuery({
    queryKey: [`workflows/${workflowId}/jobs`],
    queryFn: () =>
      fetchData<WithData<Step[]>>(
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
      fetchData<WithData<CommandOutput[]>>(
        `${API_URL}/workflows/${workflowId}/jobs/${jobId}/steps/${stepId}/output`
      ),
    refetchInterval: 500,
    enabled,
  });
};
