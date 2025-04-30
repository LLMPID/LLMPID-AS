import api from "@/utils/axiosInstance";

export const fetchExternalSystems = async () => {
  try {
    const res = await api.get(`/system/external`);
    return res.data;
  } catch (error) {
    throw new Error("Failed to fetch external systems");
  }
};

export const addExternalSystem = async (system_name: string) => {
  try {
    const res = await api.post(`/system/external`, { system_name });
    return res.data;
  } catch (error) {
    throw new Error("Failed to add external system");
  }
};

export const deleteExternalSystem = async (systemName: string) => {
  try {
    const res = await api.delete(
      `/system/external/${encodeURIComponent(systemName)}`
    );
    return res.data;
  } catch (error) {
    throw new Error("Failed to delete external system");
  }
};
