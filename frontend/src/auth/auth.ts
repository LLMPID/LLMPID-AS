import api from '@/utils/axiosInstance';

export const login = async (loginData: { username: string; password: string }) => {
  const response = await api.post('/user/auth/login', loginData);
  return response.data;
};