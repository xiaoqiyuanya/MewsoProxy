import axios, { AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { MessagePlugin } from 'tdesign-vue-next';
import router from '@/router';
import { getToken, clearToken } from '@/utils/token';

interface ApiResp<T = unknown> {
  code: number;
  message: string;
  data: T;
}

const CODE_OK = 0;
const CODE_NOT_LOGIN = 20001;
const CODE_TOKEN_EXPIRED = 20003;

const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
  withCredentials: true,
});

service.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

let refreshing = false;
let queue: Array<(token: string) => void> = [];

async function refreshToken(): Promise<string> {
  const { data } = await axios.post<ApiResp<{ access_token: string }>>(
    `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/auth/refresh`,
    {},
    { withCredentials: true },
  );
  if (data.code !== CODE_OK) {
    throw new Error(data.message);
  }
  return data.data.access_token;
}

function onLogout() {
  clearToken();
  router.replace('/login');
}

service.interceptors.response.use(
  async (response: AxiosResponse<ApiResp>): Promise<any> => {
    const res = response.data;
    if (res.code === CODE_OK) {
      return res.data;
    }
    if (res.code === CODE_NOT_LOGIN || res.code === CODE_TOKEN_EXPIRED) {
      if (!refreshing) {
        refreshing = true;
        try {
          const token = await refreshToken();
          localStorage.setItem('mewsoproxy_access_token', token);
          queue.forEach((cb) => cb(token));
          queue = [];
          return service(response.config as AxiosRequestConfig);
        } catch (e) {
          queue.forEach((cb) => cb(''));
          queue = [];
          onLogout();
          return Promise.reject(e);
        } finally {
          refreshing = false;
        }
      }
      return new Promise((resolve, reject) => {
        queue.push((token: string) => {
          if (token) {
            resolve(service(response.config as AxiosRequestConfig));
          } else {
            reject(new Error(res.message));
          }
        });
      });
    }
    MessagePlugin.error(res.message || '请求失败');
    return Promise.reject(new Error(res.message));
  },
  (error) => {
    MessagePlugin.error(error?.response?.data?.message || error.message || '网络异常');
    return Promise.reject(error);
  },
);

export default service;
