import { useCallback, useState } from 'react';

import { Args, ICallData } from '@massalabs/massa-web3';
import { bytesToStr } from '@massalabs/web3-utils';

import { usePrepareScCall } from './usePrepareScCall';
import { contracts } from '@/utils/const';

export function useMNS() {
  const [address, setAddress] = useState<string>('');
  const [mns, setMns] = useState<string>('');

  const { client, account } = usePrepareScCall();

  const accountAddress = account?.address() ?? '';

  const reverseResolveData = {
    targetAddress: contracts.buildnet.MNSContract,
    targetFunction: 'dnsReverseResolve',
    parameter: new Args().addString(accountAddress).serialize(),
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

    const resolveData = {
      targetAddress: contracts.buildnet.MNSContract,
      targetFunction: 'dnsResolve',
      parameter: new Args().addString(domain).serialize(),
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