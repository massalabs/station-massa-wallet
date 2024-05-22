import { useCallback, useState } from 'react';

import { Args } from '@massalabs/massa-web3';
import { bytesToStr } from '@massalabs/web3-utils';

import { useMassaWeb3Store } from '@/store/store';
import { contracts } from '@/utils/const';

export function useMNS() {
  const { defaultClient: client, isMainnet } = useMassaWeb3Store();

  const [targetMnsAddress, setTargetMnsAddress] = useState<string>('');
  const [domainNameList, setDomainNameList] = useState<string[]>([]);

  const targetContractAddress = isMainnet
    ? contracts.mainnet.mnsContract
    : contracts.buildnet.mnsContract;

  const reverseResolveDns = useCallback(
    async (targetAddress: string) => {
      if (!client || !targetAddress) return;
      try {
        const result = await client.smartContracts().readSmartContract({
          targetAddress: targetContractAddress,
          targetFunction: 'dnsReverseResolve',
          parameter: new Args().addString(targetAddress).serialize(),
        });
        const domains = bytesToStr(result.returnValue).split(',');

        setDomainNameList(domains);
        return domains;
      } catch (e) {
        setDomainNameList([]);
        console.error(e);
      }
    },
    [client, targetContractAddress],
  );

  const resolveDns = useCallback(
    async (domain: string): Promise<string | undefined> => {
      try {
        if (!client) return;
        const result = await client.smartContracts().readSmartContract({
          targetAddress: targetContractAddress,
          targetFunction: 'dnsResolve',
          parameter: new Args().addString(domain).serialize(),
        });

        const targetAddress = bytesToStr(result.returnValue);
        setTargetMnsAddress(bytesToStr(result.returnValue));
        return targetAddress;
      } catch (e) {
        console.error(e);
        setTargetMnsAddress('');
        return '';
      }
    },
    [client, targetContractAddress],
  );

  function resetTargetMnsAddress() {
    setTargetMnsAddress('');
  }

  const resetDomainList = useCallback(() => {
    setDomainNameList([]);
  }, []);

  return {
    domainNameList,
    resolveDns,
    targetMnsAddress,
    reverseResolveDns,
    resetTargetMnsAddress,
    resetDomainList,
  };
}
