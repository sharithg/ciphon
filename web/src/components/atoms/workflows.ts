import { atom } from "jotai";
import { TJobs, TWorkflowRunInfo } from "@/types/api";

export const workflows = atom<TWorkflowRunInfo[] | null>(null);
export const jobs = atom<TJobs[] | null>(null);
