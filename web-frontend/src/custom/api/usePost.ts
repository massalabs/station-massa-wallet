// STYLES

// EXTERNALS
import axios from 'axios';
import { UseMutationResult, useMutation } from '@tanstack/react-query';

// LOCALS

function usePost<T>(
  resource: string,
  nickname: string,
): UseMutationResult<T, unknown, T, unknown> {
  // All of our routes using POST method as nickname in the url
  const url = `${import.meta.env.VITE_BASE_API}/${resource}/${nickname}`;

  return useMutation<T, unknown, T, unknown>({
    mutationKey: [resource],
    mutationFn: async (payload) => {
      const { data } = await axios.post<T>(url, payload);

      return data;
    },
  });
}

export default usePost;
