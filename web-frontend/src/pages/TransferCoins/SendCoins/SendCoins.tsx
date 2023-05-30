import { useState } from 'react';
import usePost from '../../../../custom/api/usePost';
import { SendTransactionObject } from '../../../../models/AccountModel';
import { Button } from '@massalabs/react-ui-kit';
import { useLocation } from 'react-router-dom';

export function SendCoins() {
  const { state } = useLocation();
  const nickname = state.nickname;
  const { mutate, isSuccess } = usePost<SendTransactionObject>(
    'accounts',
    nickname,
    'transfer',
  );
  if (isSuccess) {
    console.log('Transfer success');
  }

  const [amount, setAmount] = useState<string>('64');
  const [fees, setFees] = useState<string>('42');
  const [recipient, setRecipient] = useState<string>('recipient');

  const handleTransfer = () => {
    const transferData: SendTransactionObject = {
      amount: amount,
      recipientAddress: recipient,
      fee: fees,
    };

    mutate(transferData as SendTransactionObject);
  };

  if (isSuccess) {
    console.log('Transfer success');
  }

  return (
    <div className="bg-primary w-full h-screen">
      <h1>Send coins</h1>

      <form className="flex flex-col w-1/2">
        <input
          type="text"
          value={amount}
          placeholder="amount"
          onChange={(e) => setAmount(e.target.value)}
        />
        <input
          type="text"
          value={fees}
          placeholder="fees"
          onChange={(e) => setFees(e.target.value)}
        />
        <input
          type="text"
          value={recipient}
          placeholder="recipient address"
          onChange={(e) => setRecipient(e.target.value)}
        />
        <Button onClick={handleTransfer}>Send</Button>
      </form>
    </div>
  );
}
