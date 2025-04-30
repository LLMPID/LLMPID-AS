import api from '@/utils/axiosInstance';

export const login = async (loginData: { username: string; password: string }) => {
  const response = await api.post('/user/auth/login', loginData);
  return response.data;
};

export const changePassword = async (passwordData: { username: string; old_password: string; new_password: string }) => {
  const response = await api.post('/user/auth/credentials/change', passwordData);
  return response.data;
};

