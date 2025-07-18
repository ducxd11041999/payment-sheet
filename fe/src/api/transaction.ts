import axiosInstance from "./axios";

export const getTransactions = async (month: string) => {
    const response = await axiosInstance.get(`/blocks/${month}/transactions`);
    return response.data;
};

export const addTransaction = async (month: string, data: any) => {
    const response = await axiosInstance.post(`/blocks/${month}/transactions`, data);
    return response.data;
};
