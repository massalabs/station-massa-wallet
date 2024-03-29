// STYLES

// EXTERNALS
import { UseMutationResult, useMutation } from '@tanstack/react-query';
import axios from 'axios';

// LOCALS

export function useDelete<T>(
  resource: string,
): UseMutationResult<T, unknown, T, unknown> {
  var url = `${import.meta.env.VITE_BASE_API}/${resource}`;

  return useMutation<T, unknown, T, unknown>({
    mutationKey: [resource],
    mutationFn: async () => {
      const { data } = await axios.delete<T>(url);

      return data;
    },
  });
}
