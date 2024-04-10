import { useState, useEffect } from 'react';

import { fromMAS } from '@massalabs/massa-web3';
import { toast } from '@massalabs/react-ui-kit';
import { useNavigate, useParams } from 'react-router-dom';

import { SendConfirmation, SendConfirmationData } from './SendConfirmation';
import { SendForm } from './SendForm';
import { usePost } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { AccountObject, SendTransactionObject } from '@/models/AccountModel';
import { routeFor, maskAddress } from '@/utils';

interface SendCoinsProps {
  account: AccountObject;
  redirect: {
    amount: string;
    to: string;
  };
}

export default function SendCoins(props: SendCoinsProps) {
  const { account, redirect } = props;

  const navigate = useNavigate();
  const { nickname } = useParams();

  const [submit, setSubmit] = useState<boolean>(false);
  const [data, setData] = useState<SendConfirmationData>({
    amount: '',
    fee: '',
    recipientAddress: '',
  });

  const { mutate, isSuccess, isLoading, error } =
    usePost<SendTransactionObject>(`accounts/${nickname}/transfer`);

  useEffect(() => {
    if (error) {
      toast.error(Intl.t('errors.send-coins.sent'));
    } else if (isSuccess) {
      let { amount, recipientAddress } = data;
      toast.success(
        Intl.t('success.send-coins.sent', {
          amount,
          recipient: maskAddress(recipientAddress),
        }),
      );

      navigate(routeFor(`${nickname}/home`));
    }
  }, [isSuccess, error, data, nickname, navigate]);

  function handleSubmit(data: SendConfirmationData) {
    setData(data);

    setSubmit(true);
  }

  function handleConfirm(confirmed: boolean) {
    if (!confirmed) {
      setSubmit(false);
    } else {
      mutate({
        fee: fromMAS(data.fee).toString(),
        recipientAddress: data.recipientAddress,
        amount: fromMAS(data.amount).toString(),
      });
    }
  }

  return (
    <div className="mt-5" data-testid="send-coins">
      {submit ? (
        <SendConfirmation
          data={data}
          handleConfirm={handleConfirm}
          isLoading={isLoading}
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
