import { atomWithStorage } from 'jotai/utils';
export const authTokenAtom = atomWithStorage<string | null>('auth_token', null);