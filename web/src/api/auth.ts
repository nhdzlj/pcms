import request from "./request";

export interface LoginParams {
  username: string;
  password: string;
}

export interface RegisterParams {
  username: string;
  password: string;
  email?: string;
}

export interface UserInfo {
  id: number;
  username: string;
  email: string;
  avatar: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResult {
  token: string;
  user: UserInfo;
}

export function login(params: LoginParams) {
  return request.post<any, AuthResult>("/auth/login", params);
}

export function register(params: RegisterParams) {
  return request.post<any, AuthResult>("/auth/register", params);
}

export function getMe() {
  return request.get<any, UserInfo>("/auth/me");
}
