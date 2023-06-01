import useResource from '../../../custom/api/useResource';
import { AccountObject } from '../../../models/AccountModel';
import Intl from '../../../i18n/i18n';
import { useLocation } from 'react-router-dom';
import { parseForm } from '../../../utils/parseForm';
import { useState, ChangeEvent, FormEvent } from 'react';
import {
  formatStandard,
  reverseFormatStandard,
  checkRecipientFormat,
} from '../../../utils/MassaFormating';
import { SendForm } from './SendForm';
import { SendConfirmation } from './SendConfirmation';

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
  let account = accounts.find((account) => account.nickname === nickname);
  const unformattedBalance = account?.candidateBalance ?? '0';
  const balance = parseInt(unformattedBalance) / 10 ** 9;

  const [amount, setAmount] = useState<string>('');
  const [recipient, setRecipient] = useState('');
  const [valid, setValid] = useState<boolean>(false);
  const [error, setError] = useState<IErrorObject | null>(null);

  // TODO : implement address check to see if it is real

  function validate(e: any) {
    const formObject = parseForm(e);
    const { amount, recipient } = formObject;

    const amountNum = reverseFormatStandard(amount);

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
    if (!checkRecipientFormat(recipient)) {
      setError({ recipient: Intl.t('errors.send.invalid-recipient') });
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
    else setValid(!valid);
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
        className="mr-3.5 hover:cursor-pointer"
      >
        {percentage === 100 ? 'Max' : `${percentage}%`}
      </li>
    );
  }

  function handleChange(e: ChangeEvent<HTMLInputElement>): void {
    e.preventDefault();
    setAmount(e.target.value);
  }

  const sendArgs = {
    amount: amount,
    formattedBalance: formattedBalance,
    recipient: recipient,
    error: error,
    setRecipient,
    handleSubmit,
    handleChange,
    SendPercentage,
  };
  const confirmArgs = {
    amount: amount,
    recipient: recipient,
    valid: valid,
    setValid,
  };

  return (
    <div className="mt-3.5">
      {valid ? (
        <SendConfirmation {...confirmArgs} />
      ) : (
        <SendForm {...sendArgs} />
      )}
    </div>
  );
}

export default Send;
