// STYLES

// EXTERNALS
import { useQuery, UseQueryResult } from '@tanstack/react-query';
import axios, { AxiosResponse } from 'axios';

// LOCALS

export function useResource<T>(resource: string): UseQueryResult<T, undefined> {
  const url = `${import.meta.env.VITE_BASE_API}/${resource}`;

  return useQuery<T, undefined>({
    queryKey: ['', url],
    queryFn: async () => {
      const { data } = await axios.get<T, AxiosResponse<T>>(url);

      return data;
    },
  });
}
