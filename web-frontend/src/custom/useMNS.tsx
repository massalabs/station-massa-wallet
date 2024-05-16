import { useCallback, useState } from 'react';

import { Args, Client } from '@massalabs/massa-web3';
import { bytesToStr } from '@massalabs/web3-utils';

import { usePrepareScCall } from './usePrepareScCall';
import { contracts } from '@/utils/const';

export function useMNS(client: Client | undefined) {
  const [targetMnsAddress, setTargetMnsAddress] = useState<string>('');
  const [domainNameList, setDomainNameList] = useState<string[]>([]);

  const { isMainnet } = usePrepareScCall();

  const targetContractAddress = isMainnet
    ? contracts.mainnet.mnsContract
    : contracts.buildnet.mnsContract;

  const reverseResolveDns = useCallback(
    async (targetAddress = '') => {
      if (!client) return;

      const result = await client.smartContracts().readSmartContract({
        targetAddress: targetContractAddress,
        targetFunction: 'dnsReverseResolve',
        parameter: new Args().addString(targetAddress).serialize(),
      });
      const domains = bytesToStr(result.returnValue).split(',');

      setDomainNameList(domains);
      return domains;
    },
    [client, targetContractAddress],
  );

  const resolveDns = useCallback(
    async (domain: string) => {
      if (!client) return;
      const result = await client.smartContracts().readSmartContract({
        targetAddress: targetContractAddress,
        targetFunction: 'dnsResolve',
        parameter: new Args().addString(domain).serialize(),
      });
      const targetAddress = bytesToStr(result.returnValue);
      setTargetMnsAddress(bytesToStr(result.returnValue));
      return targetAddress;
    },
    [client, targetContractAddress],
  );

  function resetTargetMnsAddress() {
    setTargetMnsAddress('');
  }

  return {
    domainNameList,
    resolveDns,
    targetMnsAddress,
    reverseResolveDns,
    resetTargetMnsAddress,
  };
}
