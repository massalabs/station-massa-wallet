import { useEffect, useState } from 'react';

import { CHAIN_ID, Provider } from '@massalabs/massa-web3';
import { getWallet, WalletName } from '@massalabs/wallet-provider';
import { useParams } from 'react-router-dom';

export function useProvider() {
  const { nickname } = useParams();
  const [provider, setProvider] = useState<Provider | undefined>();
  const [isMainnet, setIsMainnet] = useState<boolean>();

  useEffect(() => {
    getWallet(WalletName.MassaWallet).then((wallet) => {
      if (!wallet) {
        return;
      }

      wallet.networkInfos().then((network) => {
        const isMainnet = network.chainId === CHAIN_ID.Mainnet;
        setIsMainnet(isMainnet);
      });

      if (nickname) {
        wallet.accounts().then((accounts) => {
          const provider = accounts.find((a) => a.accountName === nickname);
          setProvider(provider);
        });
      }

      wallet.listenNetworkChanges((network) => {
        const isMainnet = network.chainId === CHAIN_ID.Mainnet;
        setIsMainnet(isMainnet);
      });
    });
  }, [nickname, setProvider, setIsMainnet]);

  return { provider, isMainnet };
}
