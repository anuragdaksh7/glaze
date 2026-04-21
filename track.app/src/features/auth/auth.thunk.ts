import { createAsyncThunk } from "@reduxjs/toolkit";
import type { AuthUser, User } from "./types";
import { loginUser, me, signupUser, verifyGoogleOAuth } from "./auth.api";

export const signUpUser = createAsyncThunk<{ id: string }, { name: string; email: string; password: string }>(
  "auth/signUpUser",
  async ({ name, email, password }, { rejectWithValue }) => {
    try {
      return await signupUser(name, email, password);
    } catch (error) {
      return rejectWithValue("Failed to sign up user");
    }
  }
)

export const signInUser = createAsyncThunk<AuthUser, { email: string; password: string }>(
  "auth/signInUser",
  async ({ email, password }, { rejectWithValue }) => {
    try {
      return await loginUser(email, password);
    } catch (error) {
      return rejectWithValue("Failed to sign in user");
    }
  }
)

export const verifyGoogleCode = createAsyncThunk<AuthUser, string>(
  "auth/verifyGoogleCode",
  async (code: string, { rejectWithValue }) => {
    try {
      return await verifyGoogleOAuth(code);
    } catch (error) {
      return rejectWithValue("Failed to verify Google code");
    }
  }
)

export const getMe = createAsyncThunk<User>(
  "auth/getMe",
  async (_, { rejectWithValue }) => {
    try {
      return await me();
    } catch (error) {
      return rejectWithValue("Failed to fetch user data");
    }
  }
)
