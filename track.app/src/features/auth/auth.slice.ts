import { createSlice } from "@reduxjs/toolkit";
import type { AuthUser, User } from "./types";
import { SLICE_NAMES } from "../../shared/constants/enums";
import { getMe, signInUser, verifyGoogleCode } from "./auth.thunk";

interface UserSlice {
  user: User | null;
  token: string | null;
}

const data: AuthUser | null = localStorage.getItem(SLICE_NAMES.USER) ? JSON.parse(localStorage.getItem(SLICE_NAMES.USER)!) : null;
const user: User | null = {
  id: data?.id!,
  name: data?.name!,
  email: data?.email!,
  role: data?.role!,
  profilePicture: data?.profilePicture!,
};

const initialState: UserSlice = data ? {
  user: user,
  token: data?.token || null,
} : {
  user: null,
  token: null,
}

const userSlice = createSlice({
  name: SLICE_NAMES.USER,
  initialState,
  reducers: {},
  extraReducers: builder => {
    builder
      .addCase(signInUser.fulfilled, (state, action) => {
        localStorage.setItem(SLICE_NAMES.USER, JSON.stringify(action.payload));
        state.user = {
          id: action.payload.id,
          name: action.payload.name,
          email: action.payload.email,
          role: action.payload.role,
          profilePicture: action.payload.profilePicture,
        };
        state.token = action.payload.token;
      })
      .addCase(signInUser.rejected, (state) => {
        localStorage.removeItem(SLICE_NAMES.USER);
        state.user = null;
        state.token = null;
      })
      .addCase(signInUser.pending, (state) => {
        state.user = null;
        state.token = null;
      })

      .addCase(verifyGoogleCode.fulfilled, (state, action) => {
        localStorage.setItem(SLICE_NAMES.USER, JSON.stringify(action.payload));
        state.user = {
          id: action.payload.id,
          name: action.payload.name,
          email: action.payload.email,
          role: action.payload.role,
          profilePicture: action.payload.profilePicture,
        };
        state.token = action.payload.token;
      })
      .addCase(verifyGoogleCode.rejected, (state) => {
        localStorage.removeItem(SLICE_NAMES.USER);
        state.user = null;
        state.token = null;
      })
      .addCase(verifyGoogleCode.pending, (state) => {
        state.user = null;
        state.token = null;
      })

      .addCase(getMe.fulfilled, (state, action) => {
        state.user = {
          id: action.payload.id,
          name: action.payload.name,
          email: action.payload.email,
          role: action.payload.role,
          profilePicture: action.payload.profilePicture,
        };
      })
      .addCase(getMe.rejected, (state) => {
        state.user = null;
        state.token = null;
      })
  }
})

export default userSlice.reducer;
