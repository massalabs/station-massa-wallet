import { useState, useEffect } from 'react';

import { Mas } from '@massalabs/massa-web3';
import { toast } from '@massalabs/react-ui-kit';
import { maskAddress } from '@massalabs/react-ui-kit/src/lib/massa-react/utils';
import { useNavigate, useParams } from 'react-router-dom';

import { MAS } from '@/const/assets/assets';
import { usePost } from '@/custom/api';
import { useFTTransfer } from '@/custom/smart-contract/useFTTransfer';
import Intl from '@/i18n/i18n';
import { AccountObject, SendTransactionObject } from '@/models/AccountModel';
import {
  SendConfirmation,
  SendConfirmationData,
} from '@/pages/TransferCoins/SendCoins/SendConfirmation';
import { SendForm } from '@/pages/TransferCoins/SendCoins/SendForm';
import { Redirect } from '@/pages/TransferCoins/TransferCoins';
import { useAppStore } from '@/store/store';
import { routeFor } from '@/utils';

interface SendCoinsProps {
  account?: AccountObject;
  redirect: Redirect;
}

export default function SendCoins(props: SendCoinsProps) {
  const { account, redirect } = props;

  const navigate = useNavigate();
  const { nickname } = useParams();
  const { setDisableSwitchAccount } = useAppStore();

  const [submit, setSubmit] = useState<boolean>(false);
  const [data, setData] = useState<SendConfirmationData>();
  const [errorToastId, setErrorToastId] = useState<string | null>(null);

  const {
    mutate: transferMAS,
    isSuccess: transferMASSuccess,
    isLoading: transferMASLoading,
    error: transferMASError,
  } = usePost<SendTransactionObject>(`accounts/${nickname}/transfer`);
  const {
    transfer: transferFT,
    isOpPending: transferFTPending,
    isPending: transferFTLoading,
  } = useFTTransfer();

  useEffect(() => {
    if (transferMASError && !errorToastId) {
      setErrorToastId(toast.error(Intl.t('errors.send-coins.sent')));
    } else if (transferMASSuccess && data) {
      let { amount, recipientAddress } = data;
      toast.success(
        Intl.t('success.send-coins.sent', {
          amount,
          recipient: maskAddress(recipientAddress),
        }),
      );

      navigate(routeFor(`${nickname}/home`));
      setDisableSwitchAccount(false);
    } else if (transferFTPending) {
      navigate(routeFor(`${nickname}/home`));
      setDisableSwitchAccount(false);
    }
  }, [
    transferMASSuccess,
    transferMASError,
    transferFTPending,
    data,
    nickname,
    navigate,
    errorToastId,
    setErrorToastId,
    setDisableSwitchAccount,
  ]);

  useEffect(() => {
    setDisableSwitchAccount(submit && !!data);
  }, [submit, data, setDisableSwitchAccount]);

  function handleSubmit(data: SendConfirmationData) {
    setData(data);

    setSubmit(true);
  }

  function handleConfirm(confirmed: boolean) {
    if (!confirmed) {
      setSubmit(false);
    } else if (data) {
      if (data.asset.symbol !== MAS && data.asset.address) {
        if (!transferFT) {
          return;
        }
        transferFT(
          data.recipientAddress,
          data.asset.address,
          data.amount,
          data.asset.decimals,
          data.fees,
        );
      } else {
        setErrorToastId(null);
        transferMAS({
          fee: Mas.fromString(data.fees).toString(),
          recipientAddress: data.recipientAddress,
          amount: Mas.fromString(data.amount).toString(),
        });
      }
    }
  }

  if (!account) return null;

  return (
    <div className="mt-5" data-testid="send-coins">
      {submit && data ? (
        <SendConfirmation
          data={data}
          handleConfirm={handleConfirm}
          isLoading={transferMASLoading || !!transferFTLoading}
        />
      ) : (
        <SendForm
          handleSubmit={handleSubmit}
          redirect={redirect}
          data={data}
          account={account}
        />
      )}
    </div>
  );
}
