import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});
export const getAllAccounts = () => api.get('accounts').then((res) => res.data);
