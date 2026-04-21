import { axios_instance } from "../../shared/lib/axios"
import type { AuthUser, User } from "./types";

export const signupUser = async (name: string, email: string, password: string): Promise<{ id: string }> => {
  const res = await axios_instance.post("/user/sign-up", {
    name: name,
    email: email,
    password: password,
  })
  return res.data.data;
}

export const loginUser = async (email: string, password: string): Promise<AuthUser> => {
  const res = await axios_instance.post("/user/sign-in", {
    email: email,
    password: password,
  })
  return res.data.data;
}

export const getGoogleRedirectURL = async (): Promise<{ redirect_url: string }> => {
  const res = await axios_instance.get("/oauth/google/redirect-uri");
  console.log(res.data.data)
  return res.data.data;
}

export const verifyGoogleOAuth = async (code: string): Promise<AuthUser> => {
  const res = await axios_instance.post("/oauth/google/verify?code=" + encodeURIComponent(code));
  return res.data.data;
}

export const me = async (): Promise<User> => {
  const res = await axios_instance.get("/user/me");
  return res.data.data;
}