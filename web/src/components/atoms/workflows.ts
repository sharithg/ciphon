import { atom } from "jotai";
import { Job, WorklfowRun } from "../../hooks/react-query/use-workflows";

export const workflows = atom<WorklfowRun[] | null>(null);
export const jobs = atom<Job[] | null>(null);
