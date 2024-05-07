import { MAINNET_CHAIN_ID } from '@massalabs/massa-web3';
import { providers } from '@massalabs/wallet-provider';
import { create } from 'zustand';

export interface NetworkStoreState {
  isMainnet: boolean;
  setIsMainnet: (newValue: boolean) => void;
  initNetworkStore: () => void;
}

export const useNetworkStore = create<NetworkStoreState>(
  (set: (params: Partial<NetworkStoreState>) => void) => ({
    isMainnet: false,
    setIsMainnet: (setMainnet: boolean) => {
      set({ isMainnet: setMainnet });
    },

    initNetworkStore: async () => {
      console.log('initModeStore');
      try {
        const massaProvider = (await providers())?.find(
          (p) => p.name() === 'MASSASTATION',
        );
        if (!massaProvider) {
          throw new Error('FATAL: Massa provider not found');
        }

        const chainId = await massaProvider.getChainId();

        if (chainId === MAINNET_CHAIN_ID) {
          set({ isMainnet: true });
        } else {
          set({ isMainnet: false });
        }
      } catch (e) {
        console.error('Error initializing Massa provider', e);
      }
    },
  }),
);
