import { create } from 'zustand';

import { AppStoreState, appStore } from './appStore';

export const useAppStore = create<AppStoreState>((set) => ({
  ...appStore(set),
}));
