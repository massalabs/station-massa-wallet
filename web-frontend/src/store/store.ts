import { create } from 'zustand';

import { AppStoreState, appStore } from './appStore';
import { MassaWeb3StoreState, massaWeb3Store } from './massaWeb3Store';

export const useAppStore = create<AppStoreState>((set) => ({
  ...appStore(set),
}));

export const useMassaWeb3Store = create<MassaWeb3StoreState>((set, get) => ({
  ...massaWeb3Store(set, get),
}));

function initStores() {
  useMassaWeb3Store.getState().initNetwork();
}

initStores();
