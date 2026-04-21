import type { RootState } from "../../app/store";

export const selectUser = (state: RootState) => state.user.user;
export const selectUserToken = (state: RootState) => state.user.token;
export const selectIsAuthenticated = (state: RootState) => !!state.user.token;
