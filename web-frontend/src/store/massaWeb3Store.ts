import {
  BUILDNET,
  CHAIN_ID,
  Client,
  ClientFactory,
  DefaultProviderUrls,
  MAINNET,
  MAINNET_CHAIN_ID,
} from '@massalabs/massa-web3';
import { providers } from '@massalabs/wallet-provider';
import { MASSA_STATION_PROVIDER_NAME } from '@massalabs/wallet-provider/dist/esm/massaStation/MassaStationProvider';

export interface MassaWeb3StoreState {
  defaultClient: Client | undefined;
  setDefaultClient: (client: Client | undefined) => void;

  isMainnet: boolean | undefined;
  setMainnet: (isMainnet: boolean) => void;

  initMassaDefaultClient: () => void;
  initNetwork: () => void;
}

export const massaWeb3Store = (
  set: (params: Partial<MassaWeb3StoreState>) => void,
  get: () => MassaWeb3StoreState,
) => ({
  defaultClient: undefined,
  setDefaultClient: (publicClient: Client | undefined) => {
    set({ defaultClient: publicClient });
  },

  isMainnet: undefined,
  setMainnet: (isMainnet: boolean) => {
    set({ isMainnet });
  },

  initMassaDefaultClient: async () => {
    const isMainnet = get().isMainnet;
    const massaClient = await ClientFactory.createDefaultClient(
      isMainnet ? DefaultProviderUrls.MAINNET : DefaultProviderUrls.BUILDNET,
      isMainnet ? CHAIN_ID[MAINNET] : CHAIN_ID[BUILDNET],
    );

    set({ defaultClient: massaClient });
  },

  initNetwork: async () => {
    try {
      const massaProvider = (await providers())?.find(
        (p) => p.name() === MASSA_STATION_PROVIDER_NAME,
      );
      if (!massaProvider) {
        throw new Error('FATAL: Massa provider not found');
      }

      const chainId = await massaProvider.getChainId();

      massaProvider.listenNetworkChanges((network) => {
        set({ isMainnet: network === MAINNET.toLowerCase() });
        get().initMassaDefaultClient();
      });

      set({ isMainnet: chainId === MAINNET_CHAIN_ID });
    } catch (e) {
      console.error('Error initializing Massa provider', e);
    }
  },
});
