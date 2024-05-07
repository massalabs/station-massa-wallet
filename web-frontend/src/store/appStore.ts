export interface AppStoreState {
  disableSwitchAccount: boolean;
  setDisableSwitchAccount: (newValue: boolean) => void;
}

export const appStore = (set: (params: Partial<AppStoreState>) => void) => ({
  disableSwitchAccount: false,
  setDisableSwitchAccount: (disableSwitchAccount: boolean) => {
    set({ disableSwitchAccount });
  },
});
