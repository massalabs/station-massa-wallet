// STYLES

// EXTERNALS
import { create } from 'zustand';

// LOCALS
import accountStore from './accountStore';
import configStore, { ConfigStoreState } from './configStore';

export const useAppStore = create(() => ({
  ...accountStore(),
}));

export const useConfigStore = create<ConfigStoreState>((...obj) => ({
  ...configStore(...obj),
}));
