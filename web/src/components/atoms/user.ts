import { atomWithStorage } from "jotai/utils";
import { TGetUserByIdRow } from "../../types/api";

export const userAtom = atomWithStorage<TGetUserByIdRow | null>("user", null);
