import axios from "axios";
import type { AxiosInstance, AxiosResponse } from "axios";
import { ElMessage } from "element-plus";
import router from "@/router";

const request: AxiosInstance = axios.create({
  baseURL: "/api/v1",
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求拦截器 —— 附加 Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器 —— 自动解包 {code, message, data} 返回 data
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data;
    if (res.code !== 0) {
      ElMessage.error(res.message || "请求失败");
      return Promise.reject(new Error(res.message));
    }
    return res.data;
  },
  (error) => {
    if (error.response) {
      const status = error.response.status;
      const msg = error.response.data?.message || "请求失败";

      if (status === 401) {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        router.push("/login");
        ElMessage.error("登录已过期，请重新登录");
      } else {
        ElMessage.error(msg);
      }
    } else {
      ElMessage.error("网络异常，请检查网络连接");
    }
    return Promise.reject(error);
  }
);

export default request;
