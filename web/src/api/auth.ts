import request from '@/utils/request';

export interface TokenDTO {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface UserDTO {
  id: number;
  email: string;
  balance: number;
  commission_balance: number;
  is_admin: boolean;
  is_staff: boolean;
  banned: boolean;
  plan_id?: number;
  group_id?: number;
  expired_at: string;
  token: string;
  uuid: string;
  transfer_enable: number;
  used_traffic: number;
  created_at: string;
}

export interface AuthResp {
  token: TokenDTO;
  user: UserDTO;
}

export interface LoginReq {
  email: string;
  password: string;
}

export interface RegisterReq {
  email: string;
  password: string;
  invite_code?: string;
}

export function login(data: LoginReq) {
  return request.post<AuthResp, AuthResp>('/auth/login', data);
}

export function register(data: RegisterReq) {
  return request.post<AuthResp, AuthResp>('/auth/register', data);
}

export function logout() {
  return request.post<null, null>('/user/logout');
}

export function getMe() {
  return request.get<UserDTO, UserDTO>('/user/me');
}

export interface SubscribeResp {
  token: string;
  url: string;
}

export function getSubscribe() {
  return request.get<SubscribeResp, SubscribeResp>('/user/subscribe');
}
