import axios from 'axios';
import { accountType } from './types';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});

export const getAllAccounts = (): Promise<accountType[]> =>
  api.get('accounts').then((res) => res.data);

export const getAccount = (nickname: string): Promise<accountType> =>
  api.get(`accounts/${nickname}`).then((res) => res.data);

export const createAccount = (nickname: string): Promise<accountType> =>
  api.post(`accounts/${nickname}`).then((res) => res.data);

export const importAccount = (): Promise<accountType> =>
  api.put(`accounts`).then((res) => res.data);
