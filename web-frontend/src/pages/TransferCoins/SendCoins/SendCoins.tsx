import { useState, useEffect } from 'react';

import { fromMAS } from '@massalabs/massa-web3';
import { toast } from '@massalabs/react-ui-kit';
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
import { routeFor, maskAddress } from '@/utils';

interface SendCoinsProps {
  account?: AccountObject;
  redirect: Redirect;
}

export default function SendCoins(props: SendCoinsProps) {
  const { account, redirect } = props;

  const navigate = useNavigate();
  const { nickname } = useParams();

  const [submit, setSubmit] = useState<boolean>(false);
  const [data, setData] = useState<SendConfirmationData>();
  const [errorToastId, setErrorToastId] = useState<string | null>(null);

  const {
    mutate: transferMAS,
    isSuccess: transferMASSuccess,
    isLoading: transferMASLoading,
    error: transferMASError,
  } = usePost<SendTransactionObject>(`accounts/${nickname}/transfer`);
  const { transfer: transferFT, isPending: transferFTLoading } = useFTTransfer(
    nickname || '',
  );

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
    }
  }, [
    transferMASSuccess,
    transferMASError,
    data,
    nickname,
    navigate,
    errorToastId,
    setErrorToastId,
  ]);

  function handleSubmit(data: SendConfirmationData) {
    setData(data);

    setSubmit(true);
  }

  function handleConfirm(confirmed: boolean) {
    if (!confirmed) {
      setSubmit(false);
    } else if (data) {
      if (data.asset.symbol !== MAS) {
        transferFT(
          data.recipientAddress,
          data.asset.address,
          data.amount,
          data.asset.decimals,
        );
      } else {
        setErrorToastId(null);
        transferMAS({
          fee: fromMAS(data.fees).toString(),
          recipientAddress: data.recipientAddress,
          amount: fromMAS(data.amount).toString(),
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
          isLoading={transferMASLoading || transferFTLoading}
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
