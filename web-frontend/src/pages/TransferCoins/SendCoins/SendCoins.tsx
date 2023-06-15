import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { routeFor } from '../../../utils';
import { usePost } from '../../../custom/api';
import { SendForm } from './SendForm';
import { SendConfirmation } from './SendConfirmation';
import { SendTransactionObject } from '../../../models/AccountModel';
import { toNanoMASS } from '../../../utils/massaFormat';

function SendCoins({ ...props }) {
  const { account, redirect } = props;

  const navigate = useNavigate();
  const { nickname } = useParams();

  const [submit, setSubmit] = useState<boolean>(false);
  const [data, setData] = useState<object>();
  const [payloadData, setPayloadData] = useState<object>();

  const { mutate, isSuccess, isLoading, error } =
    usePost<SendTransactionObject>(`accounts/${nickname}/transfer`);

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    } else if (isSuccess) {
      navigate(routeFor(`${nickname}/home`));
    }
  }, [isSuccess]);

  function handleSubmit({ ...data }) {
    setData(data);

    setPayloadData({
      fee: data.fees,
      recipientAddress: data.recipient,
      amount: toNanoMASS(data.amount).toString(),
    });

    setSubmit(true);
  }

  function handleConfirm(confirmed: boolean) {
    if (!confirmed) {
      setSubmit(false);
    } else {
      mutate(payloadData as SendTransactionObject);
    }
  }

  return (
    <div className="mt-5">
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

export default SendCoins;
