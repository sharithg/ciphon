import { atomWithStorage } from "jotai/utils";
import { User } from "../../hooks/user-auth";

export const userAtom = atomWithStorage<User | null>("user", null);
