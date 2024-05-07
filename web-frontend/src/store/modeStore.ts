import { Client, ClientFactory, MAINNET_CHAIN_ID } from '@massalabs/massa-web3';
import { providers } from '@massalabs/wallet-provider';
import { create } from 'zustand';

export interface ModeStoreState {
  isMainnet: boolean;
  setIsMainnet: (newValue: boolean) => void;
  nickname: string;
  setNickname: (nickname: string) => void;
  accountClient?: Client;
  setAccountClient: (client: Client) => void;
  initModeStore: (nickname: string | undefined) => void;
}

export const useModeStore = create<ModeStoreState>(
  (set: (params: Partial<ModeStoreState>) => void) => ({
    isMainnet: false,
    setIsMainnet: (setMainnet: boolean) => {
      set({ isMainnet: setMainnet });
    },
    nickname: '',
    setNickname: (nickname: string) => {
      set({ nickname });
    },
    accountClient: undefined,
    setAccountClient: (client: Client) => {
      set({ accountClient: client });
    },

    initModeStore: async (nickname) => {
      if (!nickname) {
        throw new Error('FATAL: No nickname provided');
      }

      set({ nickname: nickname });
      try {
        const massaProvider = (await providers())?.find(
          (p) => p.name() === 'MASSASTATION',
        );
        if (!massaProvider) {
          throw new Error('FATAL: Massa provider not found');
        }

        const account = (await massaProvider.accounts()).find(
          (a) => a.name() === nickname,
        );
        if (!account) {
          return;
        }

        const client = await ClientFactory.fromWalletProvider(
          massaProvider,
          account,
        );

        set({ accountClient: client });

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
