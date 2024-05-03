import { create } from 'zustand';

interface AppStoreState {
  disableSwitchAccount: boolean;
  setDisableSwitchAccount: (newValue: boolean) => void;
}

export const useAppStore = create<AppStoreState>(
  (set: (params: Partial<AppStoreState>) => void) => ({
    disableSwitchAccount: false,
    setDisableSwitchAccount: (disableSwitchAccount: boolean) => {
      set({ disableSwitchAccount });
    },
  }),
);
