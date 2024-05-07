import { useEffect, useState } from 'react';

import {
  Client,
  EOperationStatus,
  ICallData,
  MAX_GAS_CALL,
  Args,
  ClientFactory,
  MAINNET_CHAIN_ID,
  strToBytes,
  STORAGE_BYTE_COST,
} from '@massalabs/massa-web3';
import { ToastContent, parseAmount, toast } from '@massalabs/react-ui-kit';
import { providers } from '@massalabs/wallet-provider';

import { logSmartContractEvents } from './massa-utils';
import { OperationToast } from '@/components/OperationToast';
import Intl from '@/i18n/i18n';

async function prepareFTTransfer(nickname: string) {
  const massaProvider = (await providers())?.find(
    (p) => p.name() === 'MASSASTATION',
  );
  if (!massaProvider) {
    console.error('FATAL: Massa provider not found');
    return;
  }

  const account = (await massaProvider.accounts()).find(
    (a) => a.name() === nickname,
  );
  if (!account) {
    return;
  }

  return {
    client: await ClientFactory.fromWalletProvider(massaProvider, account),
    chainId: await massaProvider.getChainId(),
  };
}

export function useFTTransfer(nickname: string) {
  const [client, setClient] = useState<Client>();
  const [chainId, setChainId] = useState<bigint>();
  const isMainnet = chainId === MAINNET_CHAIN_ID;
  useEffect(() => {
    prepareFTTransfer(nickname).then((result) => {
      setClient(result?.client);
      setChainId(result?.chainId);
    });
  }, [nickname, client, setClient]);
  const {
    opId,
    isPending,
    isOpPending,
    isSuccess,
    isError,
    callSmartContract,
  } = useWriteSmartContract(client, isMainnet);

  const transfer = (
    recipient: string,
    tokenAddress: string,
    amount: string,
    decimals: number,
  ) => {
    if (!client) {
      throw new Error('Massa client not found');
    }

    const rawAmount = parseAmount(amount, decimals);
    const args = new Args().addString(recipient).addU256(rawAmount);

    estimateCoinsCost(client, tokenAddress, recipient).then((coins) => {
      callSmartContract(
        'transfer',
        tokenAddress,
        args.serialize(),
        {
          pending: Intl.t('send-coins.steps.ft-transfer-pending'),
          success: Intl.t('send-coins.steps.ft-transfer-success'),
          error: Intl.t('send-coins.steps.ft-transfer-failed'),
        },
        coins,
      );
    });
  };

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

interface ToasterMessage {
  pending: string;
  success: string;
  error: string;
  timeout?: string;
}

function minBigInt(a: bigint, b: bigint) {
  return a < b ? a : b;
}

function useWriteSmartContract(client?: Client, isMainnet?: boolean) {
  const [isPending, setIsPending] = useState(false);
  const [isOpPending, setIsOpPending] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const [isError, setIsError] = useState(false);
  const [opId, setOpId] = useState<string | undefined>(undefined);

  function callSmartContract(
    targetFunction: string,
    targetAddress: string,
    parameter: number[],
    messages: ToasterMessage,
    coins = BigInt(0),
  ) {
    if (!client) {
      throw new Error('Massa client not found');
    }
    if (isOpPending) {
      throw new Error('Operation is already pending');
    }
    setIsSuccess(false);
    setIsError(false);
    setIsOpPending(false);
    setIsPending(true);
    let operationId: string | undefined;
    let toastId: string | undefined;

    const callData = {
      targetAddress,
      targetFunction,
      parameter,
      coins,
    } as ICallData;

    client
      .smartContracts()
      .readSmartContract({
        ...callData,
        callerAddress: client.wallet().getBaseAccount()?.address(),
      })
      .then((response) => {
        const gasCost = BigInt(response.info.gas_cost);
        return minBigInt(gasCost + (gasCost * 20n) / 100n, MAX_GAS_CALL);
      })
      .then((maxGas: bigint) => {
        callData.maxGas = maxGas;
        return client.smartContracts().callSmartContract(callData);
      })
      .then((opId) => {
        operationId = opId;
        setOpId(operationId);
        setIsOpPending(true);
        toastId = toast.loading(
          (t) => (
            <ToastContent t={t}>
              <OperationToast
                isMainnet={isMainnet}
                title={messages.pending}
                operationId={operationId}
              />
            </ToastContent>
          ),
          {
            duration: Infinity,
          },
        );
        return client
          .smartContracts()
          .awaitMultipleRequiredOperationStatus(operationId, [
            EOperationStatus.SPECULATIVE_ERROR,
            EOperationStatus.SPECULATIVE_SUCCESS,
          ]);
      })
      .then((status: EOperationStatus) => {
        if (status !== EOperationStatus.SPECULATIVE_SUCCESS) {
          throw new Error('Operation failed', { cause: { status } });
        }
        setIsSuccess(true);
        setIsOpPending(false);
        setIsPending(false);
        toast.dismiss(toastId);
        toast.success((t) => (
          <ToastContent t={t}>
            <OperationToast
              isMainnet={isMainnet}
              title={messages.success}
              operationId={operationId}
            />
          </ToastContent>
        ));
      })
      .catch((error) => {
        console.error(error);
        toast.dismiss(toastId);
        setIsError(true);
        setIsOpPending(false);
        setIsPending(false);

        if (!operationId) {
          console.error('Operation ID not found');
          toast.error((t) => (
            <ToastContent t={t}>
              <OperationToast isMainnet={isMainnet} title={messages.error} />
            </ToastContent>
          ));
          return;
        }

        if (error.cause?.status === EOperationStatus.SPECULATIVE_ERROR) {
          toast.error((t) => (
            <ToastContent t={t}>
              <OperationToast
                isMainnet={isMainnet}
                title={messages.error}
                operationId={operationId}
              />
            </ToastContent>
          ));
          logSmartContractEvents(client, operationId);
        } else {
          toast.error((t) => (
            <ToastContent t={t}>
              <OperationToast
                isMainnet={isMainnet}
                title={
                  messages.timeout || Intl.t('send-coins.steps.failed-timeout')
                }
                operationId={operationId}
              />
            </ToastContent>
          ));
        }
      });
  }

  return {
    opId,
    isOpPending,
    isPending,
    isSuccess,
    isError,
    callSmartContract,
  };
}
async function estimateCoinsCost(
  client: Client,
  tokenAddress: string,
  recipient: string,
): Promise<bigint> {
  const addrInfo = await client.publicApi().getAddresses([tokenAddress]);
  const allKeys = addrInfo[0].candidate_datastore_keys;
  const key = balanceKey(recipient);
  const foundKey = allKeys.find((k) => k === key);

  if (foundKey) {
    return 0n;
  }

  const storage =
    4n + // space of a key/value in the datastore
    BigInt(key.length) + // key length
    32n; // length of the value of the balance

  return STORAGE_BYTE_COST * storage;
}

function balanceKey(address: string): number[] {
  const BALANCE_KEY_PREFIX = 'BALANCE';

  return Array.from(strToBytes(BALANCE_KEY_PREFIX + address));
}
