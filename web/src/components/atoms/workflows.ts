import { atom } from "jotai";

type TItem = {
  name: string;
  id: string;
};

export const selectedWorkflowAtom = atom<TItem | null>(null);
export const selectedJobAtom = atom<TItem | null>(null);
