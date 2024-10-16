import { atom } from "jotai";
import { TGetJobsByWorkflowIdRow, TGetWorkflowRunsRow } from "@/types/api";

export const workflows = atom<TGetWorkflowRunsRow[] | null>(null);
export const jobs = atom<TGetJobsByWorkflowIdRow[] | null>(null);
