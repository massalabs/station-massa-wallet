import { useCallback, useState } from 'react';

import { Args, bytesToStr, SmartContract } from '@massalabs/massa-web3';

import { useProvider } from './useProvider';
import { contracts } from '@/utils/const';

export function useMNS() {
  const [targetMnsAddress, setTargetMnsAddress] = useState<string>('');
  const [domainNameList, setDomainNameList] = useState<string[]>([]);

  const { provider, isMainnet } = useProvider();

  const targetContractAddress = isMainnet
    ? contracts.mainnet.mnsContract
    : contracts.buildnet.mnsContract;

  const reverseResolveDns = useCallback(
    async (targetAddress: string) => {
      if (!provider) {
        return;
      }
      try {
        const mnsContract = new SmartContract(provider, targetContractAddress);
        const res = await mnsContract.read(
          'dnsReverseResolve',
          new Args().addString(targetAddress),
        );
        const domains = bytesToStr(res.value)
          .split(',')
          .filter((d) => !!d.length);

        setDomainNameList(domains);
        return domains;
      } catch (e) {
        setDomainNameList([]);
        console.error(e);
      }
    },
    [provider, targetContractAddress],
  );

  const resolveDns = useCallback(
    async (domain: string): Promise<string | undefined> => {
      if (!provider) {
        return;
      }
      try {
        const mnsContract = new SmartContract(provider, targetContractAddress);

        const res = await mnsContract.read(
          'dnsResolve',
          new Args().addString(domain),
        );
        const targetAddress = bytesToStr(res.value);
        setTargetMnsAddress(targetAddress);
        return targetAddress;
      } catch (e) {
        console.error(e);
        setTargetMnsAddress('');
      }
    },
    [provider, targetContractAddress],
  );

  const resetTargetMnsAddress = useCallback(() => {
    setTargetMnsAddress('');
  }, []);
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
