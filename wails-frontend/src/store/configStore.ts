// STYLES

// EXTERNALS
import { persist } from 'zustand/middleware';

// LOCALS

export interface ConfigStoreState {
  timeoutId: number | undefined;
}

const configStore = persist<ConfigStoreState>(
  () => ({
    timeoutId: undefined
  }),
  {
    name: 'config-store',
  },
);

export default configStore;
