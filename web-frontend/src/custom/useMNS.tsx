import { useCallback, useState } from 'react';

import { Args, ICallData, MAINNET_CHAIN_ID } from '@massalabs/massa-web3';
import { bytesToStr } from '@massalabs/web3-utils';

import { usePrepareScCall } from './usePrepareScCall';
import { contracts, networks } from '@/utils/const';

export function useMNS() {
  const [targetMnsAddress, setTargetMnsAddress] = useState<string>('');
  const [mns, setMns] = useState<string>('');

  const { client, account, chainId } = usePrepareScCall();

  const currentNetwork =
    chainId === MAINNET_CHAIN_ID ? networks.mainnet : networks.buildnet;

  const reverseResolveDns = useCallback(
    async (address = '') => {
      const targetAddress = address || account?.address();
      if (!targetAddress) return;
      const reverseResolveCallData = {
        targetAddress: contracts[currentNetwork].mnsContract,
        targetFunction: 'dnsReverseResolve',
        parameter: new Args().addString(targetAddress).serialize(),
      } as ICallData;

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
    },
    [client],
  );

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
          setTargetMnsAddress(bytesToStr(result.returnValue));
        });
    } catch (e) {
      setTargetMnsAddress('');
      console.error(e);
    }
  }

  function resetTargetMnsAddress() {
    setTargetMnsAddress('');
  }

  return {
    mns,
    resolveDns,
    targetMnsAddress,
    reverseResolveDns,
    resetTargetMnsAddress,
  };
}
