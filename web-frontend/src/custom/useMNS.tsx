import { useCallback, useState } from 'react';

import { Args, ICallData } from '@massalabs/massa-web3';
import { bytesToStr } from '@massalabs/web3-utils';

import { usePrepareScCall } from './usePrepareScCall';
import { contracts } from '@/utils/const';

export function useMNS() {
  const [targetMnsAddress, setTargetMnsAddress] = useState<string>('');
  const [mns, setMns] = useState<string>('');

  const { client, account, isMainnet } = usePrepareScCall();

  const targetContractAddress = isMainnet
    ? contracts.mainnet.mnsContract
    : contracts.buildnet.mnsContract;

  const reverseResolveDns = useCallback(
    async (address = '') => {
      const targetAddress = address || account?.address();
      if (!targetAddress) return;

      const reverseResolveCallData = {
        targetAddress: targetContractAddress,
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
        console.error('Reverse DNS resolution failed:', e);
      }
    },
    [client, account?.address, targetContractAddress],
  );

  const resolveDns = useCallback(
    async (domain: string) => {
      if (!client) return;
      const resolveCallData = {
        targetAddress: targetContractAddress,
        targetFunction: 'dnsResolve',
        parameter: new Args().addString(domain).serialize(),
      } as ICallData;

      try {
        const result = await client
          .smartContracts()
          .readSmartContract(resolveCallData);
        setTargetMnsAddress(bytesToStr(result.returnValue));
      } catch (e) {
        setTargetMnsAddress('');
        console.error('DNS resolution failed:', e);
      }
    },
    [client, targetContractAddress],
  );

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
