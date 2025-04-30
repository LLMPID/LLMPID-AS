import axios from 'axios';
import { authTokenAtom } from '@/atoms/authAtom';
import { getDefaultStore } from 'jotai';

const api = axios.create({
  baseURL: '/api', 
});

const store = getDefaultStore();

api.interceptors.request.use(
  (config) => {
    const token = store.get(authTokenAtom);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        store.set(authTokenAtom, null); 
      }
      return Promise.reject(error);
    }
  );

export default api;