import { useEffect, useState } from 'react';

import { Args, Client, ICallData } from '@massalabs/massa-web3';
import { IAccount } from '@massalabs/wallet-provider';
import { bytesToStr } from '@massalabs/web3-utils';

import { prepareSCCall } from '@/utils/prepareSCCall';

export function useMNS(nickname: string | undefined) {
  const [client, setClient] = useState<Client>();
  const [account, setAccount] = useState<IAccount>();
  const [mns, setMns] = useState<string>('');

  if (!nickname) {
    throw new Error('Nickname not found');
  }

  const MNSTargetAddress =
    'AS1CpitsdLu4dtbQrqAzhThygL2ytGyacFED1ogr2HsxZxfNy8qQ';

  const accountAddress = account?.address() ?? '';

  const reverseResolveArgs = new Args().addString(accountAddress);

  const reverseResolveData = {
    targetAddress: MNSTargetAddress,
    targetFunction: 'dnsReverseResolve',
    parameter: reverseResolveArgs.serialize(),
  } as ICallData;

  useEffect(() => {
    prepareSCCall(nickname).then((result) => {
      setClient(result?.client);
      setAccount(result?.account);
    });
  }, [nickname, client, setClient]);

  async function reverseResolveDns() {
    if (!client) return;
    try {
      await client
        .smartContracts()
        .readSmartContract(reverseResolveData)
        .then((result) => {
          setMns(bytesToStr(result.returnValue));
        });
    } catch (e) {
      console.error(e);
    }
  }

  return { reverseResolveDns, mns };
}
