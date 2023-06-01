import { Balance, Button, Input } from '@massalabs/react-ui-kit';
import useResource from '../../../custom/api/useResource';
import { AccountObject } from '../../../models/AccountModel';
import Intl from '../../../i18n/i18n';
import { useLocation } from 'react-router-dom';
import { parseForm } from '../../../utils/parseForm';
import { useState, ChangeEvent, FormEvent } from 'react';
import {
  formatStandard,
  reverseFormatStandard,
} from '../../../utils/MassaFormating';

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

    if (amountNum > balance) {
      setError({ amount: Intl.t('errors.send.amount-to-high') });
      return false;
    }
    if (Number.isNaN(amountNum)) {
      setError({ amount: Intl.t('errors.send.invalid-amount') });
      return false;
    }
    if (!recipient) {
      setError({ recipient: Intl.t('errors.send.no-recipient') });
      return false;
    }
    setError(null);
    return true;
  }

  // This function is going to be used for the backend
  // Some parts are commented to avoid a build errors
  // because the amount and recipient aren't used yet
  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!validate(e)) return;
    // const formObject = parseForm(e);
    // const { amount, recipient } = formObject;
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

  function handleChange(e: ChangeEvent<HTMLInputElement>): void {
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
