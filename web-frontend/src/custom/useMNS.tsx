import { useCallback, useEffect, useState } from 'react';

import { Args, Client, ICallData } from '@massalabs/massa-web3';
import { IAccount } from '@massalabs/wallet-provider';
import { bytesToStr } from '@massalabs/web3-utils';

import { prepareSCCall } from '@/utils/prepareSCCall';

export function useMNS(nickname: string | undefined) {
  const [client, setClient] = useState<Client>();
  const [account, setAccount] = useState<IAccount>();
  const [address, setAddress] = useState<string>('');
  const [mns, setMns] = useState<string>('');
  if (!nickname) {
    throw new Error('Nickname not found');
  }

  useEffect(() => {
    prepareSCCall(nickname).then((result) => {
      setClient(result?.client);
      setAccount(result?.account);
    });
  }, [nickname, setClient]);

  const MNSTargetAddress =
    'AS1CpitsdLu4dtbQrqAzhThygL2ytGyacFED1ogr2HsxZxfNy8qQ';

  const accountAddress = account?.address() ?? '';

  const reverseResolveArgs = new Args().addString(accountAddress);

  const reverseResolveData = {
    targetAddress: MNSTargetAddress,
    targetFunction: 'dnsReverseResolve',
    parameter: reverseResolveArgs.serialize(),
  } as ICallData;

  const reverseResolveDns = useCallback(async () => {
    if (!client) return;
    try {
      const result = await client
        .smartContracts()
        .readSmartContract(reverseResolveData);
      setMns(bytesToStr(result.returnValue));
    } catch (e) {
      setMns('');
      console.error(e);
    }
  }, [client, reverseResolveData]);

  async function resolveDns(domain: string): Promise<string | undefined> {
    if (!client) return;
    const resolveArgs = new Args().addString(domain);
    const resolveData = {
      targetAddress: MNSTargetAddress,
      targetFunction: 'dnsResolve',
      parameter: resolveArgs.serialize(),
    } as ICallData;
    try {
      await client
        .smartContracts()
        .readSmartContract(resolveData)
        .then((result) => {
          setAddress(bytesToStr(result.returnValue));
        });
    } catch (e) {
      setAddress('');
      console.error(e);
    }
  }

  return { mns, resolveDns, address, reverseResolveDns };
}
