import axiosInstance from "./axios";

export const login = async (username: string, password: string) => {
    const response = await axiosInstance.post("/login", { username, password });
    return response.data;
};

export const register = async (username: string, password: string) => {
    const response = await axiosInstance.post("/register", { username, password });
    return response.data;
};
