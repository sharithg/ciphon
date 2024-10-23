import { TEdge, TGetJobsByWorkflowIdRow } from "./api";

export type TJobsResponse = {
  jobs: TGetJobsByWorkflowIdRow[];
  edges: TEdge[];
};
