// STYLES

// EXTERNALS
import axios, { AxiosResponse } from 'axios';
import { UseMutationResult, useMutation } from '@tanstack/react-query';

// LOCALS

function usePost<T>(
  resource: string,
): UseMutationResult<T, unknown, T, unknown> {
  const url = `${import.meta.env.VITE_BASE_APP}/${resource}`;

  return useMutation<T, unknown, T, unknown>({
    mutationKey: [resource],
    mutationFn: async (payload) => {
      const { data } = await axios.post<T, AxiosResponse<T>>(url, payload);

      return data;
    },
  });
}

export default usePost;
