import axios from "axios";

export const API_URL = process.env.REACT_APP_API_URL || "http://localhost:3000";

const api = axios.create({
  baseURL: API_URL,
});

// Interceptor: tự động thêm Authorization header cho mỗi request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Transactions
export const getTransactions = (month: string) =>
  api.get(`/blocks/${month}/transactions`);

export const getMembers = (month: string) =>
  api.get(`/blocks/${month}/members`);

export const addTransaction = (
  month: string,
  data: {
    description: string;
    amount: number;
    payer: string;
    ratios: Record<string, number>;
  }
) => api.post(`/blocks/${month}/transactions`, data);

export const deleteTransaction = (id: string) =>
  api.delete(`/transactions/${id}`);

// Blocks
export const getBlocks = () => api.get("/blocks");

export const createBlock = (month: string, members: string[]) => {
  const memberArray = members.map((name) => ({ name: name.trim() }));
  return api.post("/blocks", { month, members: memberArray });
};

export const toggleLock = (month: string, locked: boolean) =>
  api.post(`/blocks/${month}/${locked ? "unlock" : "lock"}`, {});

// Auth
export const login = async (username: string, password: string) => {
  const res = await fetch(`${API_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.message || "Sai thông tin đăng nhập");
  }

  localStorage.setItem("token", data.token);
  localStorage.setItem("username", username);
  return data.token;
};
