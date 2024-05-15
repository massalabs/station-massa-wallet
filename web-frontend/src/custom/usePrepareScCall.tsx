import { useEffect, useState } from 'react';

import { Client, MAINNET_CHAIN_ID } from '@massalabs/massa-web3';
import { IAccount } from '@massalabs/wallet-provider';
import { useParams } from 'react-router-dom';

import { prepareSCCall } from '@/utils/prepareSCCall';

export function usePrepareScCall() {
  const { nickname } = useParams();
  const [client, setClient] = useState<Client>();
  const [chainId, setChainId] = useState<bigint>();
  const [account, setAccount] = useState<IAccount>();

  useEffect(() => {
    if (!nickname) {
      throw new Error('Nickname not found');
    }
    prepareSCCall(nickname).then((result) => {
      setClient(result?.client);
      setChainId(result?.chainId);
      setAccount(result?.account);
    });
  }, [nickname, setClient]);

  const isMainnet = chainId === MAINNET_CHAIN_ID;

  return { client, chainId, account, isMainnet };
}
