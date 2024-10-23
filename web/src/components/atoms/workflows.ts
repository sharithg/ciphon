import { TGetWorkflowRunsRow } from "@/types/api";
import { TJobsResponse } from "../../types";
import { atomWithStorage } from "jotai/utils";

export const workflows = atomWithStorage<TGetWorkflowRunsRow[] | null>(
  "workflows",
  null
);
export const jobs = atomWithStorage<TJobsResponse | null>("jobs", null);
