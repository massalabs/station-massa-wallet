import { useCallback } from 'react';

import { Args, Mas, StorageCost } from '@massalabs/massa-web3';
import { parseAmount, useWriteSmartContract } from '@massalabs/react-ui-kit';

import { useProvider } from '../useProvider';
import Intl from '@/i18n/i18n';

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
      // eslint-disable-next-line new-cap
      const coins = await StorageCost.MRC20BalanceCreationCost(
        provider,
        tokenAddress,
        recipient,
      );
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
