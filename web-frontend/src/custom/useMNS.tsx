import { useCallback, useState } from 'react';

import { Args, ICallData, MAINNET_CHAIN_ID } from '@massalabs/massa-web3';
import { bytesToStr } from '@massalabs/web3-utils';

import { usePrepareScCall } from './usePrepareScCall';
import { contracts, networks } from '@/utils/const';

export function useMNS() {
  const [address, setAddress] = useState<string>('');
  const [mns, setMns] = useState<string>('');

  const { client, account, chainId } = usePrepareScCall();

  const accountAddress = account?.address() ?? '';

  const currentNetwork =
    chainId === MAINNET_CHAIN_ID ? networks.mainnet : networks.buildnet;

  const reverseResolveCallData = {
    targetAddress: contracts[currentNetwork].mnsContract,
    targetFunction: 'dnsReverseResolve',
    parameter: new Args().addString(accountAddress).serialize(),
  } as ICallData;

  const reverseResolveDns = useCallback(async () => {
    if (!client) return;
    try {
      const result = await client
        .smartContracts()
        .readSmartContract(reverseResolveCallData);
      setMns(bytesToStr(result.returnValue));
    } catch (e) {
      setMns('');
      console.error(e);
    }
  }, [client, reverseResolveCallData]);

  async function resolveDns(domain: string): Promise<string | undefined> {
    if (!client) return;

    const resolveCallData = {
      targetAddress: contracts[currentNetwork].mnsContract,
      targetFunction: 'dnsResolve',
      parameter: new Args().addString(domain).serialize(),
    } as ICallData;
    try {
      await client
        .smartContracts()
        .readSmartContract(resolveCallData)
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
