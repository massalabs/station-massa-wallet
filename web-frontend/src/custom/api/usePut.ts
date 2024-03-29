// STYLES

// EXTERNALS
import { UseMutationResult, useMutation } from '@tanstack/react-query';
import axios, { AxiosResponse } from 'axios';

// LOCALS

export function usePut<T>(
  resource: string,
): UseMutationResult<T, unknown, T, unknown> {
  const url = `${import.meta.env.VITE_BASE_API}/${resource}`;

  return useMutation<T, unknown, T, unknown>({
    mutationKey: [resource],
    mutationFn: async (payload) => {
      const { data } = await axios.put<T, AxiosResponse<T>>(url, payload);

      return data;
    },
  });
}
