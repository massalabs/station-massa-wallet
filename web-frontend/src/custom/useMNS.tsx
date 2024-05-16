import { useCallback, useState } from 'react';

import { Args, Client, ICallData } from '@massalabs/massa-web3';
import { bytesToStr } from '@massalabs/web3-utils';

import { usePrepareScCall } from './usePrepareScCall';
import { contracts } from '@/utils/const';

export function useMNS(client: Client | undefined) {
  const [targetMnsAddress, setTargetMnsAddress] = useState<string>('');
  const [domainName, setDomainName] = useState<string>('');

  const { isMainnet } = usePrepareScCall();

  const targetContractAddress = isMainnet
    ? contracts.mainnet.mnsContract
    : contracts.buildnet.mnsContract;

  const reverseResolveDns = useCallback(
    async (targetAddress = '') => {
      const reverseResolveCallData = {
        targetAddress: targetContractAddress,
        targetFunction: 'dnsReverseResolve',
        parameter: new Args().addString(targetAddress).serialize(),
      } as ICallData;

      try {
        if (!client) throw new Error('Client not initialized');
        const result = await client
          .smartContracts()
          .readSmartContract(reverseResolveCallData);
        setDomainName(bytesToStr(result.returnValue));
      } catch (e) {
        setDomainName('');
        console.error('Reverse DNS resolution failed:', e);
      }
    },
    [client, targetContractAddress],
  );

  const resolveDns = useCallback(
    async (domain: string) => {
      const resolveCallData = {
        targetAddress: targetContractAddress,
        targetFunction: 'dnsResolve',
        parameter: new Args().addString(domain).serialize(),
      } as ICallData;

      try {
        if (!client) throw new Error('Client not initialized');
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
    domainName,
    resolveDns,
    targetMnsAddress,
    reverseResolveDns,
    resetTargetMnsAddress,
  };
}
