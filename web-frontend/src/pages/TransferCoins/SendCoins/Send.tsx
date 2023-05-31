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
  const balance = parseInt(account!.candidateBalance);
  const [amount, setAmount] = useState<number>(0);
  const [error, setError] = useState<IErrorObject | null>(null);
  const [recipient, setRecipient] = useState('');
  const selectPercentClass = 'mr-3.5 hover:cursor-pointer';

  // refactor to switch => todo
  // Check to see if address exists in our system

  function validate(e: any) {
    const formObject = parseForm(e);
    const { amount, recipient } = formObject;

    if (!amount) {
      setError({ amount: 'Amount is required' });
      return false;
    }
    const int = parseInt(amount);
    if (Number.isNaN(int)) {
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
  function handleSubmit(e: any) {
    e.preventDefault();
    if (!validate(e)) return;

    const formObject = parseForm(e);

    const { amount, recipient } = formObject;

    // console logs to avoid build fail
    console.log(amount);
    console.log(recipient);
    console.log('submit');
  }
  function formatStandard(num: number) {
    return num.toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });
  }
  const formattedBalance = formatStandard(balance);

  function overrideAmount(pct: number) {
    setAmount(balance * pct);
  }

  function SendPercentage({ percentage }: { percentage: number }) {
    return (
      <li
        onClick={() => overrideAmount(percentage / 100)}
        className={selectPercentClass}
      >
        {percentage == 100 ? 'Max' : `${percentage}%`}
      </li>
    );
  }

  // JSX form
  function AmountSection() {
    return (
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
            // value={parseToStandard(parseInt(amount ? amount : '0'))!}
            value={formatStandard(amount)}
            name="amount"
            onChange={(e) => setAmount(parseInt(e.target.value))}
            error={error?.amount}
          />
        </div>
        <div className="flex flex-row-reverse">
          <ul className="flex flex-row">
            <SendPercentage percentage={25}></SendPercentage>
            <SendPercentage percentage={50}></SendPercentage>
            <SendPercentage percentage={75}></SendPercentage>
            <SendPercentage percentage={100}></SendPercentage>
          </ul>
        </div>
      </div>
    );
  }

  function RecipientSection() {
    return (
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
    );
  }

  return (
    <div className="mt-3.5">
      <form onSubmit={handleSubmit}>
        {/* Balance Section */}
        <div>
          <p> Account Balance</p>
          <Balance
            className="pl-0"
            amount={formattedBalance}
            equal={'$0012345.67'}
          />
        </div>
        <AmountSection />
        <RecipientSection />
        {/* Button Section */}
        <div className="flex flex-col w-full">
          <Button variant={'secondary'}>Advanced</Button>
          {/* Temp fix for buttons not taking full width*/}
          <div className="mb-3.5"></div>
          <Button type="submit"> Send </Button>
        </div>
      </form>
    </div>
  );
}

export default Send;
