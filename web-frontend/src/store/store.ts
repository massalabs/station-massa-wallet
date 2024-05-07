import { create } from 'zustand';

import { AppStoreState, appStore } from './appStore';
import { NetworkStoreState, networkStore } from './networkStore';

export const useAppStore = create<AppStoreState>((set) => ({
  ...appStore(set),
}));

export const useNetworkStore = create<NetworkStoreState>((set) => ({
  ...networkStore(set),
}));

function initStores() {
  useNetworkStore.getState().initNetworkStore();
}

initStores();
