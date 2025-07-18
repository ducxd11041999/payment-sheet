import axiosInstance from "./axios";

export const createBlock = async (month: string, members: { name: string }[]) => {
    const response = await axiosInstance.post("/blocks", { month, members });
    return response.data;
};

export const getBlocks = async () => {
    const response = await axiosInstance.get("/blocks"); // Bạn cần thêm API này ở backend
    return response.data;
};
