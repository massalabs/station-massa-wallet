import { useEffect, useState } from 'react';

import { CHAIN_ID, Provider } from '@massalabs/massa-web3';
import {
  getWallet,
  MassaStationWallet,
  WalletName,
} from '@massalabs/wallet-provider';
import { useParams } from 'react-router-dom';

export function useProvider() {
  const { nickname } = useParams();
  const [provider, setProvider] = useState<Provider | undefined>();
  const [isMainnet, setIsMainnet] = useState<boolean>();
  const [wallet, setWallet] = useState<MassaStationWallet | undefined>();

  useEffect(() => {
    getWallet(WalletName.MassaWallet).then((msWallet) => {
      if (!msWallet) {
        return;
      }

      setWallet(msWallet as MassaStationWallet);

      msWallet.networkInfos().then((network) => {
        const isMainnet = network.chainId === CHAIN_ID.Mainnet;
        setIsMainnet(isMainnet);
      });

      if (nickname) {
        msWallet.accounts().then((accounts) => {
          const provider = accounts.find((a) => a.accountName === nickname);
          setProvider(provider);
        });
      }

      msWallet.listenNetworkChanges((network) => {
        const isMainnet = network.chainId === CHAIN_ID.Mainnet;
        setIsMainnet(isMainnet);
      });
    });
  }, [nickname, setProvider, setIsMainnet]);

  return { provider, isMainnet, wallet };
}
