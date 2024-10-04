import { API_URL } from "./constants";
import { useQuery } from "@tanstack/react-query";
import { fetchData } from ".";

type WorklfowRun = {
  commitSha: string;
  repoName: string;
  pipelineId: string;
  workflowId: string;
  workflowName: string;
  status: string;
  branch: string;
  createdAt: string;
  duration: number;
};

export const useGetWorkflows = () => {
  return useQuery({
    queryKey: ["nodes"],
    queryFn: () => fetchData<{ data: WorklfowRun[] }>(`${API_URL}/workflows`),
  });
};
