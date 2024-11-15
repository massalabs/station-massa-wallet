import { useCallback } from 'react';

import { Args, bytes, Mas, Provider, strToBytes } from '@massalabs/massa-web3';
import { parseAmount, useWriteSmartContract } from '@massalabs/react-ui-kit';

import { useProvider } from '../useProvider';
import Intl from '@/i18n/i18n';

const BALANCE_KEY_PREFIX = 'BALANCE';

function balanceKey(address: string): Uint8Array {
  return strToBytes(BALANCE_KEY_PREFIX + address);
}

async function estimateCoinsCost(
  provider: Provider,
  tokenAddress: string,
  recipient: string,
): Promise<bigint> {
  const allKeys = await provider.getStorageKeys(
    tokenAddress,
    BALANCE_KEY_PREFIX,
  );
  const key = balanceKey(recipient);
  const foundKey = allKeys.some((k) => {
    return JSON.stringify(k) === JSON.stringify(key);
  });

  if (foundKey) {
    return 0n;
  }

  const storage =
    4 + // space of a key/value in the datastore
    key.length + // key length
    32; // length of the value of the balance

  return bytes(storage);
}

export function useFTTransfer() {
  const { provider, isMainnet } = useProvider();
  const {
    opId,
    isPending,
    isOpPending,
    isSuccess,
    isError,
    callSmartContract,
  } = useWriteSmartContract(provider!, isMainnet);

  const transfer = useCallback(
    async (
      recipient: string,
      tokenAddress: string,
      amount: string,
      decimals: number,
      fees: string,
    ) => {
      if (!callSmartContract || !provider) {
        return;
      }
      const rawAmount = parseAmount(amount, decimals);
      const args = new Args().addString(recipient).addU256(rawAmount);
      const coins = await estimateCoinsCost(provider, tokenAddress, recipient);

      await callSmartContract(
        'transfer',
        tokenAddress,
        args.serialize(),
        {
          pending: Intl.t('send-coins.steps.ft-transfer-pending'),
          success: Intl.t('send-coins.steps.ft-transfer-success'),
          error: Intl.t('send-coins.steps.ft-transfer-failed'),
        },
        coins,
        Mas.fromString(fees),
      );
    },
    [provider, callSmartContract],
  );

  return {
    opId,
    isPending,
    isOpPending,
    isSuccess,
    isError,
    transfer,
    isMainnet,
  };
}
