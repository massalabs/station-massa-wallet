import { Balance, Button, Input } from '@massalabs/react-ui-kit';
import { FiUsers } from 'react-icons/fi';
import useResource from '../../../custom/api/useResource';
import { AccountObject } from '../../../models/AccountModel';
import { useLocation } from 'react-router-dom';
import { parseForm } from '../../../utils/parseForm';
import { useState } from 'react';

// remove padding from Balance component

interface IErrorObject {
  amount?: string;
  fees?: string;
  recipient?: string;
  error?: string;
}

function Send() {
  const { state } = useLocation();
  const nickname = state.nickname;
  const { data: accounts = [] } = useResource<AccountObject[]>('accounts');
  const account = accounts.find((account) => account.nickname === nickname);
  const balance = parseInt(account?.candidateBalance ?? '0') / 10 ** 9;
  const [amount, setAmount] = useState<string>('');
  const [error, setError] = useState<IErrorObject | null>(null);
  const [recipient, setRecipient] = useState('');

  const selectPercentClass = 'mr-3.5 hover:cursor-pointer';

  // TODO : implement address check to see if it is real

  function validate(e: any) {
    const formObject = parseForm(e);
    const { amount, recipient } = formObject;

    const amountNum = reverseFormatStandard(amount);
    console.log(amount, amountNum, balance);

    if (!amountNum && amount == '') {
      setError({ amount: 'Amount is required' });
      return false;
    }
    if (amountNum > balance) {
      setError({ amount: 'Insufficient funds on balance' });
      return false;
    }
    if (Number.isNaN(amountNum)) {
      setError({ amount: 'Please enter a valid number' });
      return false;
    }
    if (!recipient) {
      setError({ recipient: 'Recipient is required' });
      return false;
    }
    setError(null);
    return true;
  }

  // This function is going to be used for the backend
  // Some parts are commented to avoid a build errors
  // because the amount and recipient aren't used yet
  function handleSubmit(e: any) {
    e.preventDefault();
    if (!validate(e)) return;
    // const formObject = parseForm(e);
    // const { amount, recipient } = formObject;
  }

  function formatStandard(num: number, maximumFractionDigits = 2) {
    return num.toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits,
    });
  }

  function reverseFormatStandard(str: string) {
    const formattedString = str.replace(/[^0-9.-]/g, ''); // Remove non-numeric characters
    return parseFloat(formattedString);
  }

  const formattedBalance = formatStandard(balance);

  function overrideAmount(pct: number) {
    const newAmount = balance * pct;
    const newFormatedAmount = formatStandard(newAmount, 9);
    setAmount(newFormatedAmount);
  }

  function SendPercentage({ percentage }: { percentage: number }) {
    return (
      <li
        onClick={() => {
          overrideAmount(percentage / 100);
        }}
        className={selectPercentClass}
      >
        {percentage == 100 ? 'Max' : `${percentage}%`}
      </li>
    );
  }

  function handleChange(e: any) {
    e.preventDefault();
    setAmount(e.target.value);
  }

  return (
    <div className="mt-3.5">
      <form onSubmit={handleSubmit}>
        {/* Balance Section */}
        <div>
          <p> Account Balance</p>
          <Balance
            customClass="pl-0"
            amount={formattedBalance}
            equal={'0012345.67'}
          />
        </div>
        <div className="mb-3.5">
          <div className="flex flex-row justify-between w-full mb-3.5 ">
            <p>Send Token</p>
            <p>
              Available : <u>{formattedBalance}</u>
            </p>
          </div>
          <div className="mb-3.5">
            <Input
              placeholder={'Amount to send'}
              value={amount}
              name="amount"
              onChange={(e) => handleChange(e)}
              error={error?.amount}
            />
          </div>
          <div className="flex flex-row-reverse">
            <ul className="flex flex-row">
              <SendPercentage percentage={25} />
              <SendPercentage percentage={50} />
              <SendPercentage percentage={75} />
              <SendPercentage percentage={100} />
            </ul>
          </div>
        </div>
        <div>
          <div className="mb-3.5">
            <p>To :</p>
          </div>
          <div className="mb-3.5">
            <Input
              placeholder={'Recipient'}
              value={recipient}
              onChange={(e) => setRecipient(e.target.value)}
              icon={<FiUsers />}
              name="recipient"
              error={error?.recipient}
            />
          </div>
          <div className="flex flex-row-reverse mb-3.5">
            <p
              className="hover:cursor-pointer"
              onClick={() => console.log('transfer between accounts')}
            >
              <u>Transfer between my account</u>
            </p>
          </div>
        </div>
        {/* Button Section */}
        <div className="flex flex-col w-full">
          <div className="mb-3.5">
            <Button variant={'secondary'}>Advanced</Button>
          </div>
          <div>
            <Button type="submit"> Send </Button>
          </div>
        </div>
      </form>
    </div>
  );
}

export default Send;
